CREATE TABLE `Outlets` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `merchant_id` bigint(20) NOT NULL,
    `outlet_name` varchar(40) NOT NULL,
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `created_by` bigint(20) NOT NULL,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_by` bigint(20) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;
INSERT INTO Outlets Values (1, 1, 'Outlet 1', now(), 1, now(),1), (2, 2, 'Outlet 1', now(), 2, now(),2), (3, 1, 'Outlet 2', now(), 1, now(),1);