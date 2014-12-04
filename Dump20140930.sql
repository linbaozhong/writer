CREATE DATABASE  IF NOT EXISTS `writer` /*!40100 DEFAULT CHARACTER SET utf8 */;
USE `writer`;
-- MySQL dump 10.13  Distrib 5.6.17, for Win32 (x86)
--
-- Host: 127.0.0.1    Database: writer
-- ------------------------------------------------------
-- Server version	5.7.4-m14-log

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `d_accounts`
--

DROP TABLE IF EXISTS `d_accounts`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `d_accounts` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '账户id',
  `openFrom` int(11) NOT NULL DEFAULT '0' COMMENT '开放平台',
  `openId` varchar(32) NOT NULL COMMENT '开放账户id',
  `openName` varchar(45) NOT NULL COMMENT '开放名字',
  `regTime` bigint(20) NOT NULL,
  `status` int(11) NOT NULL DEFAULT '0' COMMENT '账户状态',
  PRIMARY KEY (`id`),
  KEY `openFrom` (`openFrom`,`openId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='账户表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `d_accounts`
--

LOCK TABLES `d_accounts` WRITE;
/*!40000 ALTER TABLE `d_accounts` DISABLE KEYS */;
/*!40000 ALTER TABLE `d_accounts` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `d_item_tag`
--

DROP TABLE IF EXISTS `d_item_tag`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `d_item_tag` (
  `itemId` bigint(20) NOT NULL,
  `tagId` bigint(20) NOT NULL,
  KEY `itemtag` (`itemId`,`tagId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='条目-标签';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `d_item_tag`
--

LOCK TABLES `d_item_tag` WRITE;
/*!40000 ALTER TABLE `d_item_tag` DISABLE KEYS */;
/*!40000 ALTER TABLE `d_item_tag` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `d_item_team`
--

DROP TABLE IF EXISTS `d_item_team`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `d_item_team` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `itemId` bigint(20) NOT NULL COMMENT '条目id',
  `accountId` bigint(20) NOT NULL COMMENT '账户Id',
  `role` int(11) NOT NULL DEFAULT '0' COMMENT '角色',
  PRIMARY KEY (`id`),
  KEY `itemId` (`itemId`),
  KEY `accountId` (`accountId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='条目-团队';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `d_item_team`
--

LOCK TABLES `d_item_team` WRITE;
/*!40000 ALTER TABLE `d_item_team` DISABLE KEYS */;
/*!40000 ALTER TABLE `d_item_team` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `d_items`
--

DROP TABLE IF EXISTS `d_items`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `d_items` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `parentId` bigint(20) NOT NULL DEFAULT '0' COMMENT '父Id',
  `title` varchar(100) NOT NULL DEFAULT '' COMMENT '标题',
  `intro` varchar(200) NOT NULL DEFAULT '' COMMENT '简介',
  `content` text NOT NULL COMMENT '内容',
  `position` int(11) NOT NULL DEFAULT '0' COMMENT '顺序位置',
  `status` int(11) NOT NULL DEFAULT '0' COMMENT '隐私状态',
  `deleted` int(11) NOT NULL DEFAULT '0' COMMENT '删除状态',
  `ownerId` bigint(20) NOT NULL DEFAULT '0' COMMENT '所属条目id',
  `creator` bigint(20) NOT NULL DEFAULT '0' COMMENT '创建人',
  `created` bigint(20) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updator` bigint(20) NOT NULL DEFAULT '0' COMMENT '最后修改人',
  `updated` bigint(20) NOT NULL DEFAULT '0' COMMENT '最后更新时间',
  PRIMARY KEY (`id`),
  KEY `parentId` (`parentId`),
  KEY `position` (`position`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='条目表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `d_items`
--

LOCK TABLES `d_items` WRITE;
/*!40000 ALTER TABLE `d_items` DISABLE KEYS */;
/*!40000 ALTER TABLE `d_items` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `d_tags`
--

DROP TABLE IF EXISTS `d_tags`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `d_tags` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(45) NOT NULL DEFAULT '',
  `status` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `name` (`name`),
  KEY `status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='标签';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `d_tags`
--

LOCK TABLES `d_tags` WRITE;
/*!40000 ALTER TABLE `d_tags` DISABLE KEYS */;
/*!40000 ALTER TABLE `d_tags` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping events for database 'writer'
--

--
-- Dumping routines for database 'writer'
--
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2014-09-30 16:59:51
