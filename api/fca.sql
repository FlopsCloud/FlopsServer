CREATE TABLE `user_resources` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT , 
  `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
  `resource_id` bigint unsigned NOT NULL  COMMENT '使用资源类型ID',
  `org_id` bigint unsigned NOT NULL COMMENT '组织ID',
  `resource_type` varchar(255) NOT NULL COMMENT '使用资源类型CPU,MEM,GPU,SSD,NET',
  `stat` varchar(64) NOT NULL COMMENT '状态 init,start,stop',
  `start_at` timestamp NULL DEFAULT NULL COMMENT '启动时间',
  `stop_at` timestamp NULL DEFAULT NULL COMMENT '停止时间',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `org_id` (`org_id`),
  KEY `user_id` (`user_id`),
  KEY `stat` (`stat`),
  CONSTRAINT `userresources_ibfk_1` FOREIGN KEY (`org_id`) REFERENCES `organizations` (`org_id`),
  CONSTRAINT `userresources_ibfk_2` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`),
  CONSTRAINT `userresources_ibfk_3` FOREIGN KEY (`resource_id`) REFERENCES `resources` (`resource_id`)
) 

CREATE TABLE `minute_usage` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `usage_id` bigint unsigned NOT NULL COMMENT '日使用记录ID',
  `org_id` bigint unsigned NOT NULL COMMENT '组织ID',
  `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
  `resource_id` bigint unsigned NOT NULL COMMENT '使用资源类型ID',
  `usage_date` date NOT NULL COMMENT '使用日期',
  `fee` bigint NOT NULL DEFAULT '0' COMMENT '费用',
  `discount` bigint NOT NULL DEFAULT '0' COMMENT '100 表示 10 折',
  PRIMARY KEY (`id`),
  KEY `org_id` (`org_id`),
  KEY `resource_id` (`resource_id`),
  KEY `user_id` (`user_id`),
  KEY `usage_id` (`usage_id`),
  CONSTRAINT `minute_usage_ibfk_1` FOREIGN KEY (`org_id`) REFERENCES `organizations` (`org_id`),
  CONSTRAINT `minute_usage_ibfk_2` FOREIGN KEY (`resource_id`) REFERENCES `resources` (`resource_id`),
  CONSTRAINT `minute_usage_ibfk_3` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`),
  CONSTRAINT `minute_usage_ibfk_4` FOREIGN KEY (`usage_id`) REFERENCES `daily_usage` (`usage_id`)
) 

CREATE TABLE `server_tags` (
  `id` int NOT NULL AUTO_INCREMENT,
  `server_id` bigint unsigned NOT NULL,
  `tag_id` bigint unsigned NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
)  

CREATE TABLE `tags` (
  `tag_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `tag_name` varchar(255) NOT NULL DEFAULT '' COMMENT '标签名称',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP  ,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`tag_id`)
)  


insert into tags( tag_name )  values("通用型");
insert into tags( tag_name )  values("大显存");
insert into tags( tag_name )  values("高主频CPU");
insert into tags( tag_name )  values("多核心CPU");
insert into tags( tag_name )  values("1元GPU");
insert into tags( tag_name )  values("自主可控");
insert into tags( tag_name )  values("极致性价比");

#===========================

CREATE TABLE `orgs_users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `org_id` bigint unsigned NOT NULL,
  `user_id` bigint unsigned NOT NULL,
  `role` varchar(64) NOT NULL DEFAULT 'guest' COMMENT '岗位',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_users_orgs_unique` ( `org_id`,`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci



CREATE TABLE invitation (
    id BIGINT unsigned AUTO_INCREMENT PRIMARY KEY,
    org_id BIGINT unsigned NOT NULL,
    inviter_id BIGINT unsigned NOT NULL,
    invitee_email VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL,
    status ENUM('pending', 'accepted', 'rejected') NOT NULL DEFAULT 'pending',
    invitation_token VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (org_id) REFERENCES organizations(org_id),
    FOREIGN KEY (inviter_id) REFERENCES users(user_id)
);

CREATE TABLE apply_join (
    id BIGINT unsigned AUTO_INCREMENT PRIMARY KEY,
    org_id BIGINT unsigned NOT NULL,
    user_id BIGINT unsigned NOT NULL,
    message VARCHAR(200),
    status ENUM('pending', 'approved', 'rejected') NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (org_id) REFERENCES organizations(org_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id)
);



CREATE TABLE transactions (
    trans_id bigint PRIMARY KEY AUTO_INCREMENT,
  `org_id` bigint UNSIGNED NOT NULL COMMENT '组织ID',
   `user_id` bigint  UNSIGNED NOT NULL COMMENT '用户ID',
    amount DECIMAL(10, 2) NOT NULL,
    trans_type ENUM('deposit', 'withdrawal') NOT NULL,
    description VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    KEY `idx_transactions_org_user_id` (`org_id`, `user_id`),
    FOREIGN KEY (org_id) REFERENCES organizations(org_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id)
);

-- 创建余额表12
CREATE TABLE balances (
    balance_id bigint PRIMARY KEY AUTO_INCREMENT,
    `org_id` bigint UNSIGNED NOT NULL COMMENT '组织ID',
   `user_id` bigint UNSIGNED NOT NULL COMMENT '用户ID',
    balance DECIMAL(10, 2) NOT NULL,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    KEY `idx_balances_org_user_id` (`org_id`, `user_id`),
    FOREIGN KEY (org_id) REFERENCES organizations(org_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id)
);



drop TABLE resources

CREATE TABLE resources (
    
    resource_id bigint unsigned  PRIMARY KEY AUTO_INCREMENT COMMENT '使用资源类型ID' ,
    org_id bigint unsigned NOT NULL COMMENT '组织ID', 
    resource_type VARCHAR(255) NOT NULL COMMENT '使用资源类型CPU,MEM,GPU,SSD,NET' ,
    unit_min_price bigint DEFAULT 0,
    unit_hour_price bigint DEFAULT 0,
    unit_day_price bigint DEFAULT 0,
    is_deleted int DEFAULT 0 COMMENT "是否删除",
    deleted_at timestamp NULL DEFAULT NULL COMMENT "删除时间",
   `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
   `created_by` bigint unsigned NOT NULL COMMENT "创建人ID",
    FOREIGN KEY (org_id) REFERENCES organizations(org_id),
    FOREIGN KEY (created_by) REFERENCES users(user_id)
)

drop TABLE discounts

CREATE TABLE discounts (
    discount_id bigint  unsigned  PRIMARY KEY AUTO_INCREMENT COMMENT '折扣ID',
    org_id bigint unsigned NOT NULL COMMENT '组织ID', 
    resource_id bigint unsigned NOT NULL  COMMENT '使用资源类型ID' ,
    memo varchar(255) NULL COMMENT "备注",
    startdate DATE NOT NULL COMMENT "折扣开始日期",
    enddate DATE NOT NULL COMMENT "折扣结束日期",
    discount bigint DEFAULT 0 COMMENT "100 表示 10 折",
    is_deleted int DEFAULT 0 COMMENT "是否删除",
    deleted_at timestamp NULL DEFAULT NULL COMMENT "删除时间",
   `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
   `created_by` bigint unsigned NOT NULL COMMENT "创建人ID",
    FOREIGN KEY (org_id) REFERENCES organizations(org_id),
    FOREIGN KEY (resource_id) REFERENCES resources(resource_id),
    FOREIGN KEY (created_by) REFERENCES users(user_id)
)

-- 创建每分钟消费表
CREATE TABLE minute_usage (
    id bigint unsigned PRIMARY KEY AUTO_INCREMENT,
    org_id bigint unsigned NOT NULL COMMENT '组织ID', 
    user_id bigint unsigned NOT NULL  COMMENT '用户ID',
    usage_datetime DATETIME NOT NULL COMMENT '使用日期时间',
    resource_id bigint unsigned NOT NULL  COMMENT '使用资源类型ID' ,
    usage_amount bigint DEFAULT 1  COMMENT '使用资源分钟数' ,
    FOREIGN KEY (org_id) REFERENCES organizations(org_id),
    FOREIGN KEY (resource_id) REFERENCES resources(resource_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id)
);

-- 创建全天消费表
CREATE TABLE daily_usage (
    usage_id bigint unsigned PRIMARY KEY AUTO_INCREMENT,
    org_id bigint unsigned NOT NULL  COMMENT '组织ID',
    user_id bigint unsigned NOT NULL  COMMENT '用户ID', 
    usage_date DATE NOT NULL COMMENT '使用日期' ,
    resource_id bigint unsigned NOT NULL COMMENT '使用资源类型ID' ,
    usage_min_amount bigint DEFAULT 0 COMMENT '使用资源分钟数',
    usage_hour_amount bigint DEFAULT 0 COMMENT '使用资源小时数',
    unit_hour_price bigint DEFAULT 0 COMMENT '使用资源小时价格',
    discount_id bigint unsigned   NULL COMMENT "id 空 表示不打折", 
    discount bigint DEFAULT 0  COMMENT "100 表示 10 折",
    FOREIGN KEY (org_id) REFERENCES organizations(org_id),
    FOREIGN KEY (resource_id) REFERENCES resources(resource_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id)
);


 

CREATE TABLE roles (
    role_id bigint unsigned PRIMARY KEY AUTO_INCREMENT,
    role_name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
   
);

-- User-Role mapping table (many-to-many relationship)
drop table user_roles

CREATE TABLE user_roles (
    id bigint unsigned PRIMARY KEY AUTO_INCREMENT,
    user_id bigint unsigned,
    role_id bigint unsigned,
    assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    assigned_by bigint unsigned,
   
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(role_id) ON DELETE CASCADE,
    FOREIGN KEY (assigned_by) REFERENCES users(user_id),
    INDEX idx_user_id (user_id),
    INDEX idx_role_id (role_id)
);


-- Permissions table
CREATE TABLE permissions (
    permission_id bigint unsigned PRIMARY KEY AUTO_INCREMENT,
    permission_name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE INDEX idx_permission_name (permission_name)
);

-- Role-Permission mapping table (many-to-many relationship)
drop table role_permissions


CREATE TABLE role_permissions (
    id bigint unsigned PRIMARY KEY AUTO_INCREMENT,
    role_id bigint unsigned,
    permission_id bigint unsigned,
    granted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    granted_by bigint unsigned,
    FOREIGN KEY (role_id) REFERENCES roles(role_id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permissions(permission_id) ON DELETE CASCADE,
    FOREIGN KEY (granted_by) REFERENCES users(user_id),
    INDEX idx_role_perm_role (role_id),
    INDEX idx_role_perm_permission (permission_id)
);



-- Insert basic roles
INSERT INTO roles (role_name, description) VALUES
    ('admin', 'System administrator with full access'),
    ('manager', 'Manager with elevated privileges'),
    ('user', 'Regular user with basic access');

-- Insert basic permissions
INSERT INTO permissions (permission_name, description) VALUES
    ('user.create', 'Can create new users'),
    ('user.read', 'Can view user information'),
    ('user.update', 'Can update user information'),
    ('user.delete', 'Can delete users'),
    ('role.assign', 'Can assign roles to users');

-- Assign basic permissions to roles
INSERT INTO role_permissions (role_id, permission_id) 
SELECT r.role_id, p.permission_id
FROM roles r
CROSS JOIN permissions p
WHERE r.role_name = 'admin';