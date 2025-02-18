package main

import (
	"context"
	"database/sql"
	"fca/api/internal/logic"
	"fca/api/internal/svc"
	"fca/model"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// FeeCalculator 费用计算器结构体，用于处理实例和资源的费用计算
type FeeCalculator struct {
	svcCtx *svc.ServiceContext
}

const jw = "your_token_here"

// StartFeeTimer 启动定时任务，处理费用计算和实例状态更新
// 包含三个定时器：
// 1. 每分钟计算使用费用和服务器使用情况
// 2. 每分钟更新实例状态
// 3. 每小时检查并删除过期实例
func StartFeeTimer(fc *FeeCalculator) {
	logx.Info("Start calculate usage fees")

	// 每分钟计算使用费用和服务器使用情况
	ticker := time.NewTicker(1 * time.Minute) //// ticker := time.NewTicker(10 * time.Second)
	// logx.Info("======================")
	// fc.CalculateAndInsertResourceUsage()
	logx.Info("======================")
	fc.NewCalculateAndInsertMinuteServerUsage()
	logx.Info("======================")
	fc.DeleteTimeOutRechargeOrder(context.Background())
	logx.Info("======================")

	go func() {
		for range ticker.C {
			// logx.Info("Start calculate usage fees loop")
			// if err := fc.CalculateAndInsertResourceUsage(); err != nil {
			// 	logx.Errorf("Failed to calculate usage fees: %v", err)
			// }
			if err := fc.NewCalculateAndInsertMinuteServerUsage(); err != nil {
				logx.Errorf("Failed to calculate minute server usage: %v", err)
			}

		}
	}()

	// 每分钟更新实例状态
	ticker2 := time.NewTicker(1 * time.Minute)

	go func() {
		for range ticker2.C {
			fc.UpdateInstanceStatus()
			fc.DeleteTimeOutRechargeOrder(context.Background())
		}
	}()

	// 每小时检查并删除过期实例
	EvertHour := time.NewTicker(1 * time.Hour)

	go func() {
		for range EvertHour.C {
			fc.DeleteExpiredInstances()

		}
	}()

}

func NewFeeCalculator(svcCtx *svc.ServiceContext) *FeeCalculator {
	return &FeeCalculator{
		svcCtx: svcCtx,
	}
}

// 常量定义
const (
	// 付费周期类型
	HOURLY  = "hourly"  // 按小时
	MONTHLY = "monthly" // 按月
	DAILY   = "daily"   // 按天
	YEARLY  = "yearly"  // 按年

	// 实例状态
	RUNNING = "running" // 运行中
	STOPPED = "stopped" // 已停止
	PAUSED  = "paused"  // 已暂停
	DELETED = "deleted" // 已删除
	STARTED = "started" // 已启动
)

// DeleteExpiredInstances 检查并删除过期的实例
// 1. 获取所有运行中的实例
// 2. 检查是否过期，如果过期：
//   - 删除虚拟机
//   - 更新实例状态
//   - 更新服务器GPU使用情况
//   - 停止相关资源
//   - 记录实例日志
func (fc *FeeCalculator) DeleteExpiredInstances() error {

	ctx := context.Background()

	//Get all running instances
	// conditions := fmt.Sprintf("state = '%s' and pay_cycle = '%s'", RUNNING, HOURLY)
	conditions := " 1=1 and state = 'started'"
	instances, err := fc.svcCtx.InstanceModel.FindAll(ctx, conditions)
	if err != nil {
		logx.Errorf("Failed to get running instances: %v", err)
		return err
	}

	for _, instance := range instances {

		if instance.ExpireDate.Before(time.Now()) {
			if instance.VhostStatus != STOPPED {
				// 删除虚拟机

				_, err = logic.DeleteVhost(ctx, jw, &logic.DelVhostRequest{
					Name: instance.Name,
				})
				if err != nil {
					logx.Errorf("Failed to delete vhost: %v", err)
					//return err
				}
				instance.VhostStatus = STOPPED
			}

			instance.State = DELETED
			instance.VhostStatus = STOPPED
			err = fc.svcCtx.InstanceModel.Update(ctx, instance)
			if err != nil {
				logx.Errorf("Failed to update instance: %v", err)
			}

			//update server gpu used
			Server, err := fc.svcCtx.ServerModel.FindOne(ctx, instance.ServerId)
			if err != nil {
				logx.Errorf("Failed to update instance: %v", err)
			}

			Server.GpuUsed -= int64(instance.GpuCores)

			err = fc.svcCtx.ServerModel.Update(ctx, Server)
			if err != nil {
				logx.Errorf("Failed to update instance: %v", err)

			}

			//stop resource
			err = fc.svcCtx.RunningResourcesModel.Stop(ctx, instance.InstanceId)
			if err != nil {
				logx.Errorf("Failed to stop resource: %v", err)
			}

			//insert instance log
			_, err = fc.svcCtx.InstancesLogModel.Insert(ctx, &model.InstancesLog{
				InstanceId: instance.InstanceId,
				UserId:     instance.UserId,
				Action:     "系统检测到实例过期，自动删除",
			})
			if err != nil {
				logx.Errorf("Failed to insert instance log: %v", err)
			}

		}

	}
	return nil

}

// UpdateInstanceStatus 更新实例状态
// 每分钟检查一次实例的实际运行状态，并与数据库记录同步
func (fc *FeeCalculator) UpdateInstanceStatus() error {

	// 每分钟计算一次，如果是按月计算，则计算一个月的

	ctx := context.Background()

	//Get all running instances
	// conditions := fmt.Sprintf("state = '%s' and pay_cycle = '%s'", RUNNING, HOURLY)
	conditions := " 1=1 and (state = 'started' or vhost_status = 'Running')"
	instances, err := fc.svcCtx.InstanceModel.FindAll(ctx, conditions)
	if err != nil {
		logx.Errorf("Failed to get running instances: %v", err)
		return err
	}

	vhost, err := logic.AdminListVhosts(ctx, jw)
	if err != nil {
		logx.Errorf("Failed to get vhost: %v", err)
		return err
	}

	// 遍历所有虚拟机
	for _, data := range vhost.Data {
		if data.FcbPods != nil {
			// 遍历虚拟机状态
			for _, vhoststatus := range data.FcbPods {

				// 遍历实例
				for _, instance := range instances {
					ns := "yinlf-" + strconv.Itoa(int(instance.UserId))
					logc.Info(context.Background(), instance.Name, ",", instance.UserId, ",", vhoststatus.Name, ",", ns, "=", data.Namespace, ",", instance.VhostStatus, "=", vhoststatus.Status)

					if (vhoststatus.Name == instance.Name) && (data.Namespace == ns) && (instance.VhostStatus != vhoststatus.Status) {
						instance.VhostStatus = vhoststatus.Status
						fc.svcCtx.InstanceModel.Update(ctx, instance)
					}

				}

			}

		}

	}

	// 检查是否存在未在虚拟机中运行的实例
	for _, instance := range instances {
		ns := "yinlf-" + strconv.Itoa(int(instance.UserId))
		found := false
		for _, data := range vhost.Data {
			if data.Namespace == ns {
				if data.FcbPods != nil {
					for _, vhoststatus := range data.FcbPods {
						if vhoststatus.Name == instance.Name {
							found = true
							break
						}
					}
				}
			}
		}
		// 如果实例不在虚拟机中，则更新实例状态为已停止
		if !found {
			if instance.VhostStatus != STOPPED {
				instance.VhostStatus = STOPPED
				fc.svcCtx.InstanceModel.Update(ctx, instance)
			}
		}
	}

	return nil

}

// CalculateAndInsertMinuteServerUsage 计算并记录按小时付费实例的每分钟使用费用
// 1. 获取所有运行中的按小时计费实例
// 2. 计算每分钟费用
// 3. 插入分钟使用记录
// // 4. 更新每日使用汇总
func (fc *FeeCalculator) HandleInstance(ctx context.Context, instance *model.Instances) error {

	//find latest hourly usage
	hourlyUsage, err := fc.svcCtx.HourlyUsageModel.FindLatestByInstanceId(ctx, instance.OrgId, instance.UserId, instance.InstanceId)
	if err != nil && err != sqlx.ErrNotFound {
		logx.Errorf("Failed to get latest hourly usage for instance %d: %v", instance.InstanceId, err)
		return err
	}

	var starttime time.Time
	if err == sqlx.ErrNotFound {
		starttime = instance.CreatedAt
	} else {
		if hourlyUsage.IsCharged == 0 {
			starttime = hourlyUsage.UsageDate
			fc.svcCtx.HourlyUsageModel.Delete(ctx, hourlyUsage.UsageId)

		} else {
			starttime = hourlyUsage.UsageDate.Add(time.Hour)
		}
	}

	var totalFee int64
	totalFee = 0

	starttime = starttime.Truncate(time.Second)

	now := time.Now()

	for starttime.Before(now) {

		nowHournum := DateToHourValue(now)
		nowMinnum := DateToMinValue(now)

		olddaynum := DateToInt(starttime)
		oldhournum := DateToHourValue(starttime)
		oldminnum := DateToMinValue(starttime)

		if nowHournum == oldhournum {
			//相同时间，跳出
			minuteTotal := nowMinnum - oldminnum
			fee := int64(instance.Cost) * int64(minuteTotal) / 60

			//insert hourly usage
			hourlyUsage = &model.HourlyUsage{
				UsageId:       0,
				OrgId:         instance.OrgId,
				UserId:        instance.UserId,
				InstanceId:    instance.InstanceId,
				RunresId:      0,
				Type:          0,
				UsageDate:     starttime,
				Fee:           uint64(fee),
				UnitHourPrice: int64(instance.Cost),
				Discount:      100,
				Daynum:        olddaynum,
				Hournum:       oldhournum,
				IsCharged:     0,
				MinuteBegin:   oldminnum % 100,
				MinuteEnd:     nowMinnum % 100,
				MinuteTotal:   minuteTotal,
			}
			_, err = fc.svcCtx.HourlyUsageModel.Insert(ctx, hourlyUsage)
			if err != nil {
				logx.Errorf("Failed to insert hourly usage for instance %d: %v", instance.InstanceId, err)
				continue
			}

			break
		} else {
			//不同时间，插入使用记录

			minuteTotal := 60 - (oldminnum % 100)
			fee := int64(instance.Cost) * int64(minuteTotal) / 60

			totalFee += fee

			//insert hourly usage
			hourlyUsage = &model.HourlyUsage{
				UsageId:       0,
				OrgId:         instance.OrgId,
				UserId:        instance.UserId,
				InstanceId:    instance.InstanceId,
				RunresId:      0,
				Type:          0,
				UsageDate:     starttime,
				Fee:           uint64(fee),
				UnitHourPrice: int64(instance.Cost),
				Discount:      100,
				Daynum:        olddaynum,
				Hournum:       oldhournum,
				IsCharged:     0,
				MinuteBegin:   oldminnum % 100,
				MinuteEnd:     60,
				MinuteTotal:   minuteTotal,
			}
			_, err = fc.svcCtx.HourlyUsageModel.Insert(ctx, hourlyUsage)
			if err != nil {
				logx.Errorf("Failed to insert hourly usage for instance %d: %v", instance.InstanceId, err)
				continue
			}
		}

		//remove minute & second
		starttime = starttime.Add(-time.Second * time.Duration(starttime.Second()))
		starttime = starttime.Add(-time.Minute * time.Duration(starttime.Minute()))

		starttime = starttime.Add(time.Hour)
	}

	return nil

}

func (fc *FeeCalculator) ChargeInstance(ctx context.Context, instance *model.Instances) error {
	//开启事务，获取所有is charg=0, endtime=60的记录
	//把费用记录到插入到 transaction_records 表,并设置is charged=1
	//更新instance的CNY balance

	// Start a transaction

	// hourlyUsageModel := fc.svcCtx.HourlyUsageModel.WithSession(session)
	usages, err := fc.svcCtx.HourlyUsageModel.FindUnchargedComplete(ctx, instance.OrgId, instance.UserId, instance.InstanceId)

	if err == sqlx.ErrNotFound {
		logx.Errorf("FeeTimmer ChargeInstance Failed to get uncharged hourly usages: %v", err)
		return nil
	}
	if err != nil {
		logx.Errorf("FeeTimmer ChargeInstance Failed to get uncharged hourly usages: %v", err)
		return fmt.Errorf("failed to get uncharged hourly usages: %v", err)
	}

	if len(usages) == 0 {
		return nil
	}

	var balanceValue int64

	err = fc.svcCtx.Mysql.Transact(func(session sqlx.Session) error {
		//err := sqlx.NewMysql(fc.svcCtx.Config.MySQL.DataSource).Transact(func(session sqlx.Session) error {

		// Update user's CNY balance
		hourlyUsageModel := fc.svcCtx.HourlyUsageModel.WithSession(session)
		balanceModel := fc.svcCtx.BalancesModel.WithSession(session)

		balance, err := balanceModel.FindOneByUserAndCurrency(ctx, instance.UserId, instance.OrgId, "CNY")
		if err != nil {
			return fmt.Errorf("failed to get user balance: %v", err)
		}

		transactionRecordsModel := fc.svcCtx.TransactionRecordsModel.WithSession(session)

		var totalFee uint64
		for _, usage := range usages {

			balance.Balance -= int64(usage.Fee)

			transRecord := &model.TransactionRecords{
				OrgId:     instance.OrgId,
				UserId:    instance.UserId,
				Amount:    int64(usage.Fee),
				TransType: 2,
				PayType:   "CNY",
				Detail:    "使用实例" + instance.Name + "的费用",
				OrderNo:   sql.NullString{String: "", Valid: false},
				Username:  instance.Name,
				Balance:   int64(balance.Balance),
				CreatedAt: time.Now(),
			}

			_, err = transactionRecordsModel.Insert(ctx, transRecord)
			if err != nil {
				return fmt.Errorf("failed to insert transaction record: %v", err)
			}

			usage.IsCharged = 99
			err = hourlyUsageModel.Update(ctx, usage)
			if err != nil {
				return fmt.Errorf("failed to update hourly usage: %v", err)
			}

			totalFee += usage.Fee
		}

		err = balanceModel.Update(ctx, balance)
		if err != nil {
			return fmt.Errorf("failed to update user balance: %v", err)
		}
		balanceValue = balance.Balance

		return nil
	})

	if balanceValue <= 0 {
		logx.Errorf("Balance for UserId %d is less than 0: %d", instance.UserId, balanceValue)

		instance.State = STOPPED
		instance.VhostStatus = STOPPED
		fc.svcCtx.InstanceModel.Update(context.Background(), instance)

		delresp, err := logic.DeleteVhost(ctx, jw, &logic.DelVhostRequest{
			Name: instance.Name,
		})

		if err != nil {
			logx.Errorf("Failed to delete vhost: %v", err)
		}
		if delresp.Code != 0 {
			logx.Errorf("Failed to delete vhost: %v", delresp.Info)
		}

	}

	return err
}

func (fc *FeeCalculator) NewCalculateAndInsertMinuteServerUsage() error {

	// 每分钟计算一次，如果是按月计算，则计算一个月的
	ctx := context.Background()

	// 获取所有正在运行的按小时付费实例

	conditions := fmt.Sprintf("state = '%s' and pay_cycle = '%s'", STARTED, HOURLY)
	instances, err := fc.svcCtx.InstanceModel.FindAll(ctx, conditions)
	if err != nil {
		logx.Errorf("Failed to get running instances: %v", err)
		return err
	}

	for _, instance := range instances {
		err := fc.HandleInstance(ctx, instance)
		if err != nil {
			logx.Errorf("FeeTimmer Failed to handle instance %d: %v", instance.InstanceId, err)
			continue
		}
		err = fc.ChargeInstance(ctx, instance)
		if err != nil {
			logx.Errorf("FeeTimmer Failed to charge instance %d: %v", instance.InstanceId, err)
			continue
		}
	}

	return nil
}

func (fc *FeeCalculator) CalculateAndInsertMinuteServerUsage() error {

	// 每分钟计算一次，如果是按月计算，则计算一个月的

	ctx := context.Background()

	//Get all running instances
	conditions := fmt.Sprintf("state = '%s' and pay_cycle = '%s'", STARTED, HOURLY)
	instances, err := fc.svcCtx.InstanceModel.FindAll(ctx, conditions)
	if err != nil {
		logx.Errorf("Failed to get running instances: %v", err)
		return err
	}

	now := time.Now()
	// daynum := DateToInt(now)
	// minnum := DateToMinValue(now)

	for _, instance := range instances {
		// Get resource pricing

		// Calculate per-minute fee
		minuteFee := instance.Cost / 60

		// Create minute usage record
		minuteUsage := &model.MinuteUsage{
			UsageId:       0,
			OrgId:         instance.OrgId,
			UserId:        instance.UserId,
			RunresId:      0,
			InstanceId:    instance.InstanceId,
			Type:          0,
			UsageDatetime: now,
			Fee:           uint64(minuteFee),
			Discount:      100,
			Daynum:        uint64(DateToInt(now)),
			Minnum:        uint64(DateToMinValue(now)),
		}

		// Insert minute usage
		_, err = fc.svcCtx.MinuteUsageModel.Insert(ctx, minuteUsage)
		if err != nil {
			logx.Errorf("Failed to insert minute usage for ResourceId %d: %v", instance.InstanceId, err)
			continue
		}

		// Update daily summary
		err = fc.updateDailySummary(ctx, minuteUsage, instance.Cost, now)
		if err != nil {
			logx.Errorf("Failed to update daily summary for ResourceId %d: %v", minuteUsage.RunresId, err)
		}
	}

	return nil
}

// DateToInt 将时间转换为整数格式 (YYYYMMDD)
// 例如：2024-01-01 转换为 20240101
func DateToInt(t time.Time) uint64 {
	// year, month, day := t.Date()
	// // 格式化拼接字符串，例如"20240101"这种格式
	dateStr := t.Format("20060102")
	result, _ := strconv.Atoi(dateStr)
	return uint64(result)
}

func DateToMinValue(t time.Time) uint64 {
	dateStr := t.Format("200601021504")
	result, _ := strconv.Atoi(dateStr)
	return uint64(result)
}

func DateToHourValue(t time.Time) uint64 {
	dateStr := t.Format("2006010215")
	result, _ := strconv.Atoi(dateStr)
	return uint64(result)
}

// CalculateAndInsertResourceUsage 计算并记录资源的使用费用
// 1. 获取所有正在运行的资源
// 2. 获取资源定价
// 3. 计算每分钟费用
// 4. 插入分钟使用记录
// 5. 更新每日使用汇总
func (fc *FeeCalculator) CalculateAndInsertResourceUsage() error {
	ctx := context.Background()

	runningResourcess, err := fc.svcCtx.RunningResourcesModel.FindAllRunning(ctx)
	if err != nil {
		logx.Errorf("Failed to get running runningResources: %v", err)
		return err
	}

	now := time.Now()
	daynum := DateToInt(now)
	minnum := DateToMinValue(now)

	for _, runningResources := range runningResourcess {
		// Get resource pricing
		resource, err := fc.svcCtx.ResourcesModel.FindOne(ctx, runningResources.ResourceId)
		if err != nil {
			logx.Errorf("Failed to get resource %d: %v", runningResources.ResourceId, err)
			continue
		}

		// Calculate per-minute fee
		HourlyFee := resource.UnitHourlyPrice * (100 - runningResources.HourlyDiscount) / 100

		// Create minute usage record
		minuteUsage := &model.MinuteUsage{
			UsageId:       0, // Will be set when daily record is created/updated
			OrgId:         runningResources.OrgId,
			UserId:        runningResources.UserId,
			RunresId:      runningResources.RunresId,
			InstanceId:    0,
			Type:          1, //resource
			UsageDatetime: now,
			Fee:           HourlyFee / 60,
			Discount:      int64(runningResources.HourlyDiscount), // Default to no discount (100 = 10.0)
			Daynum:        daynum,
			Minnum:        minnum,
		}

		// Insert minute usage
		_, err = fc.svcCtx.MinuteUsageModel.Insert(ctx, minuteUsage)

		if err != nil {
			if err == sqlx.ErrNotFound {
				logx.Errorf("Failed to insert minute usage for ResourceId %d: %v", runningResources.ResourceId, err)
			} else {
				logx.Errorf("%v", reflect.TypeOf(err))

				logx.Errorf("Failed to insert minute usage for ResourceId %d: %v", runningResources.ResourceId, err)
				continue
			}
		}

		// Update daily summary
		err = fc.updateDailySummary(ctx, minuteUsage, resource.UnitHourlyPrice, now)
		if err != nil {
			logx.Errorf("Failed to update daily summary for ResourceId %d: %v", runningResources.ResourceId, err)
		}
	}

	return nil
}

// updateDailySummary 更新每日使用统计
// 1. 查找当日使用记录
// 2. 如果不存在则创建新记录
// 3. 更新使用时间和费用
func (fc *FeeCalculator) updateDailySummary(ctx context.Context, minuteUsage *model.MinuteUsage, UnitHourPrice uint64, now time.Time) error {
	// Try to find existing daily record
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	// startOfHour := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location())

	// daynum := uint64(DateToInt(now))
	dailyUsage, err := fc.svcCtx.DailyUsageModel.FindByDateAndResource(
		ctx,
		uint64(minuteUsage.OrgId),
		uint64(minuteUsage.UserId),
		uint64(minuteUsage.RunresId),
		uint64(minuteUsage.Type),
		uint64(minuteUsage.Daynum),
	)

	if err == model.ErrNotFound {
		// Create new daily record
		dailyUsage = &model.DailyUsage{
			OrgId:           minuteUsage.OrgId,
			UserId:          minuteUsage.UserId,
			UsageDate:       startOfDay,
			RunresId:        minuteUsage.RunresId,
			InstanceId:      minuteUsage.InstanceId,
			UsageMinAmount:  1,
			UsageHourAmount: 0,
			Type:            minuteUsage.Type,
			UnitHourPrice:   int64(UnitHourPrice),
			Fee:             minuteUsage.Fee,
			DiscountId:      0, //不需要
			Discount:        int64(minuteUsage.Discount),
			Daynum:          minuteUsage.Daynum,
		}
		_, err = fc.svcCtx.DailyUsageModel.Insert(ctx, dailyUsage)
		if err != nil {
			logc.Errorf(ctx, "Failed to insert daily usage IN updateDailySummary for ResourceId %d: %v", dailyUsage.RunresId, err)
		}

		// 设计一个hourly 表格，记录每小时的费用，每小时更新balance

		//TODOTODO

		balance, err := fc.svcCtx.BalancesModel.FindOneByUserAndCurrency(ctx, dailyUsage.UserId, dailyUsage.OrgId, "CNY")
		if err != nil {
			logx.Errorf("Failed to get balance for UserId %d: %v", dailyUsage.UserId, err)
			return err
		}

		balance.Balance -= int64(dailyUsage.Fee)
		err = fc.svcCtx.BalancesModel.Update(ctx, balance)
		if err != nil {
			logx.Errorf("Failed to update balance for UserId %d: %v", dailyUsage.UserId, err)
			return err
		}

		if balance.Balance <= 0 {
			logx.Errorf("Balance for UserId %d is less than 0: %d", dailyUsage.UserId, balance.Balance)

		}

		return err
	} else if err != nil {
		return err
	}

	// Update existing record
	dailyUsage.UsageMinAmount++
	dailyUsage.UsageHourAmount = dailyUsage.UsageMinAmount / 60

	//TODOTODO
	dailyUsage.Fee = uint64(dailyUsage.UnitHourPrice) * uint64(dailyUsage.UsageMinAmount) / 60

	return fc.svcCtx.DailyUsageModel.Update(ctx, dailyUsage)
}

func (fc *FeeCalculator) DeleteTimeOutRechargeOrder(ctx context.Context) error {

	err := fc.svcCtx.RechargeOrderModel.DeleteTimeOut(ctx)

	return err
}
