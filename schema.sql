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
  `PASSWORD` text NOT NULL,
  `ISSYSADMIN` int(1) NOT NULL DEFAULT '1',
  `TOKEN` text NOT NULL,
  `EMAIL` text NOT NULL,
  PRIMARY KEY (`ID`),
  KEY `upass` (`USERNAME`(128),`PASSWORD`(64)),
  KEY `token` (`TOKEN`(128))
);

CREATE TABLE `logs` (
  `USERID` int(11) NOT NULL,
  `OBJECTID` text NOT NULL,
  `ACTION` text NOT NULL,
  `DETAILS` blob,
  `DATE` datetime NOT NULL
);