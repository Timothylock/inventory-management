-- This file contains all of the schemas for the tables. Import this into your blank mysql database

CREATE TABLE `items` (
  `ID` text NOT NULL,
  `NAME` text NOT NULL,
  `CATEGORY` text NOT NULL,
  `PICTURE_URL` text NOT NULL,
  `DETAILS` text NOT NULL,
  `LOCATION` text NOT NULL,
  `LAST_PERFORMED_BY` int(11) NOT NULL DEFAULT '0',
  `QUANTITY` int(11) NOT NULL DEFAULT '1',
  `STATUS` text NOT NULL,
  `DELETED` int(1) NOT NULL DEFAULT '0',
  KEY `ID_2` (`ID`(32),`DELETED`),
  FULLTEXT KEY `search` (`ID`,`NAME`,`CATEGORY`,`DETAILS`,`LOCATION`)
);

CREATE TABLE `users` (
  `ID` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `USERNAME` text NOT NULL,
  `ACTIVE` int(1) NOT NULL DEFAULT '1',
  `PASSWORD_HASHED` text NOT NULL,
  PRIMARY KEY (`ID`)
);