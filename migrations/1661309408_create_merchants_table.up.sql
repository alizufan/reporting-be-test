CREATE TABLE `Merchants` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `user_id` int(40) NOT NULL,
    `merchant_name` varchar(40) NOT NULL,
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `created_by` bigint(20) NOT NULL,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_by` bigint(20) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;
INSERT INTO Merchants VALUES (1, 1, 'merchant 1', now(), 1, now(),1), (2, 2, 'Merchant 2', now(), 2, now(),2);