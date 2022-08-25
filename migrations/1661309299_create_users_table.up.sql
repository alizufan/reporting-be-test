CREATE TABLE `Users` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `name` varchar(45) DEFAULT NULL,
    `user_name` varchar(45) DEFAULT NULL,
    `password` varchar(225) DEFAULT NULL,
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `created_by` bigint(20) NOT NULL,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_by` bigint(20) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;

INSERT INTO Users VALUES (1, 'Admin 1', 'admin1', MD5('admin1'), now(), 1, now(),1), (2, 'Admin 2', 'admin2', MD5('admin2'), now(), 2, now(),2);