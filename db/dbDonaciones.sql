CREATE DATABASE  IF NOT EXISTS `food_donation_db` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci */ /*!80016 DEFAULT ENCRYPTION='N' */;
USE `food_donation_db`;
-- MySQL dump 10.13  Distrib 8.0.33, for Win64 (x86_64)
--
-- Host: localhost    Database: food_donation_db
-- ------------------------------------------------------
-- Server version	8.0.34

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `donation_history`
--

DROP TABLE IF EXISTS `donation_history`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `donation_history` (
  `id` int NOT NULL AUTO_INCREMENT,
  `donation_id` int NOT NULL,
  `action` enum('created','reserved','completed','cancelled') NOT NULL,
  `user_id` int NOT NULL,
  `notes` text,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `donation_id` (`donation_id`),
  KEY `user_id` (`user_id`),
  CONSTRAINT `donation_history_ibfk_1` FOREIGN KEY (`donation_id`) REFERENCES `donations` (`id`) ON DELETE CASCADE,
  CONSTRAINT `donation_history_ibfk_2` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `donation_history`
--

LOCK TABLES `donation_history` WRITE;
/*!40000 ALTER TABLE `donation_history` DISABLE KEYS */;
/*!40000 ALTER TABLE `donation_history` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `donations`
--

DROP TABLE IF EXISTS `donations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `donations` (
  `id` int NOT NULL AUTO_INCREMENT,
  `donor_id` int NOT NULL,
  `title` varchar(255) NOT NULL,
  `description` text,
  `category` enum('bakery','dairy','fruits','meat','canned','prepared','sugar','fats','cereals','beverages','furniture','electronics','clothing','books','toys','appliances','tools','sports','office','other') NOT NULL,
  `quantity` int NOT NULL,
  `expiry_date` date DEFAULT NULL,
  `pickup_address` text NOT NULL,
  `pickup_latitude` decimal(10,8) DEFAULT NULL,
  `pickup_longitude` decimal(11,8) DEFAULT NULL,
  `status` enum('available','reserved','completed','expired') DEFAULT 'available',
  `reserved_by` int DEFAULT NULL,
  `reserved_at` timestamp NULL DEFAULT NULL,
  `pickup_time` time DEFAULT NULL,
  `pickup_person_name` varchar(255) DEFAULT NULL,
  `pickup_person_id` varchar(50) DEFAULT NULL,
  `verification_code` varchar(10) DEFAULT NULL,
  `completed_at` timestamp NULL DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `donor_confirmed` tinyint(1) DEFAULT '0',
  `recipient_confirmed` tinyint(1) DEFAULT '0',
  `business_confirmed` tinyint(1) DEFAULT '0',
  `business_confirmed_at` timestamp NULL DEFAULT NULL,
  `donor_confirmed_at` timestamp NULL DEFAULT NULL,
  `recipient_confirmed_at` timestamp NULL DEFAULT NULL,
  `weight` decimal(10,2) DEFAULT NULL COMMENT 'Peso en kilogramos',
  `donation_reason` text COMMENT 'Motivo de la donación',
  `contact_info` varchar(255) DEFAULT NULL COMMENT 'Información de contacto del donante',
  PRIMARY KEY (`id`),
  KEY `reserved_by` (`reserved_by`),
  KEY `idx_donations_donor` (`donor_id`),
  KEY `idx_donations_status` (`status`),
  KEY `idx_donations_category` (`category`),
  KEY `idx_donations_location` (`pickup_latitude`,`pickup_longitude`),
  KEY `idx_verification_code` (`verification_code`),
  CONSTRAINT `donations_ibfk_1` FOREIGN KEY (`donor_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `donations_ibfk_2` FOREIGN KEY (`reserved_by`) REFERENCES `users` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `donations`
--

LOCK TABLES `donations` WRITE;
/*!40000 ALTER TABLE `donations` DISABLE KEYS */;
INSERT INTO `donations` VALUES (1,1,'Pan fresco del día','Pan de diferentes tipos, horneado hoy en la mañana','bakery',20,'2025-11-16','Calle Principal 123, CDMX',19.43260000,-99.13320000,'available',NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2025-11-15 01:15:28','2025-11-15 01:15:28',0,0,0,NULL,NULL,NULL,5.50,'Excedente de producción diaria','juan@mail.com - 555-0101'),(2,1,'Verduras frescas','Zanahorias, lechugas y tomates de la huerta','fruits',15,'2025-11-17','Calle Principal 123, CDMX',19.43260000,-99.13320000,'reserved',2,'2025-11-15 01:15:28',NULL,NULL,NULL,'ABC123',NULL,'2025-11-15 01:15:28','2025-11-15 01:15:28',1,0,0,NULL,'2025-11-15 01:15:28',NULL,8.20,'Cosecha abundante','juan@mail.com - 555-0101'),(3,3,'Comida preparada','Arroz con pollo y ensalada, lista para consumir','prepared',30,'2025-11-15','Boulevard Norte 789, CDMX',19.44200000,-99.13890000,'available',NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2025-11-15 01:15:28','2025-11-15 01:15:28',0,0,1,'2025-11-15 01:15:28',NULL,NULL,12.00,'Sobrante de evento corporativo','carlos@mail.com - 555-0103'),(4,4,'Frutas variadas','Manzanas, naranjas y plátanos en buen estado','fruits',25,'2025-11-18','Calle Sur 321, CDMX',19.41950000,-99.14250000,'reserved',5,'2025-11-15 01:15:28','14:30:00','Pedro Sánchez','ID789456','XYZ789',NULL,'2025-11-15 01:15:28','2025-11-15 01:15:28',1,1,0,NULL,'2025-11-15 01:15:28','2025-11-15 01:15:28',10.30,'Limpieza de inventario','ana@mail.com - 555-0104'),(5,3,'Lácteos próximos a vencer','Leche y yogurt con fecha próxima pero en perfecto estado','dairy',40,'2025-11-16','Boulevard Norte 789, CDMX',19.44200000,-99.13890000,'completed',2,'2025-11-14 01:15:28','10:00:00','María García','ID123456','DEF456','2025-11-15 01:15:28','2025-11-15 01:15:28','2025-11-15 01:15:28',1,1,1,'2025-11-14 23:15:28','2025-11-14 22:15:28','2025-11-15 00:15:28',6.80,'Próximos a fecha de vencimiento','carlos@mail.com - 555-0103'),(6,1,'Enlatados variados','Atún, frijoles y vegetales enlatados','canned',50,'2026-05-13','Calle Principal 123, CDMX',19.43260000,-99.13320000,'available',NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2025-11-15 01:15:28','2025-11-15 01:15:28',0,0,0,NULL,NULL,NULL,15.00,'Exceso de stock','juan@mail.com - 555-0101'),(7,5,'Muebles de oficina','Sillas y escritorios en buen estado','furniture',10,NULL,'Avenida Este 654, CDMX',19.43670000,-99.12080000,'available',NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2025-11-15 01:15:28','2025-11-15 01:15:28',0,0,0,NULL,NULL,NULL,45.00,'Renovación de oficina','pedro@mail.com - 555-0105'),(8,1,'Ropa de invierno','Chamarras y suéteres en buen estado','clothing',35,NULL,'Calle Principal 123, CDMX',19.43260000,-99.13320000,'available',NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2025-11-15 01:15:28','2025-11-15 01:15:28',0,0,0,NULL,NULL,NULL,12.50,'Cambio de temporada','juan@mail.com - 555-0101'),(9,3,'Bebidas embotelladas','Agua y jugos sellados','beverages',100,'2025-12-14','Boulevard Norte 789, CDMX',19.44200000,-99.13890000,'available',NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2025-11-15 01:15:28','2025-11-15 01:15:28',0,0,0,NULL,NULL,NULL,25.00,'Sobrepedido','carlos@mail.com - 555-0103'),(10,1,'Pan del día anterior','Pan de ayer, aún en buen estado','bakery',15,'2025-11-13','Calle Principal 123, CDMX',19.43260000,-99.13320000,'expired',NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2025-11-15 01:15:28','2025-11-15 01:15:28',0,0,0,NULL,NULL,NULL,3.50,'Exceso de stock vencido','juan@mail.com - 555-0101');
/*!40000 ALTER TABLE `donations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `notifications`
--

DROP TABLE IF EXISTS `notifications`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `notifications` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL,
  `title` varchar(255) NOT NULL,
  `message` text NOT NULL,
  `type` enum('donation_created','donation_reserved','donation_completed','general') NOT NULL,
  `is_read` tinyint(1) DEFAULT '0',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_notifications_user` (`user_id`),
  KEY `idx_notifications_read` (`is_read`),
  CONSTRAINT `notifications_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `notifications`
--

LOCK TABLES `notifications` WRITE;
/*!40000 ALTER TABLE `notifications` DISABLE KEYS */;
/*!40000 ALTER TABLE `notifications` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `id` int NOT NULL AUTO_INCREMENT,
  `email` varchar(255) NOT NULL,
  `password` varchar(255) NOT NULL,
  `name` varchar(255) NOT NULL,
  `phone` varchar(20) DEFAULT NULL,
  `user_type` enum('donor','organization') NOT NULL,
  `address` text,
  `latitude` decimal(10,8) DEFAULT NULL,
  `longitude` decimal(11,8) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `donation_days` json DEFAULT NULL COMMENT 'Configuration for available donation days and times',
  `reset_token` varchar(255) DEFAULT NULL,
  `reset_token_expires` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `email` (`email`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES (1,'juan@mail.com','password123','Juan Pérez','555-0101','donor','Calle Principal 123, CDMX',19.43260000,-99.13320000,'2025-11-15 01:15:28','2025-11-15 01:15:28','[\"lunes\", \"miércoles\", \"viernes\"]',NULL,NULL),(2,'maria@mail.com','password123','María García','555-0102','organization','Avenida Central 456, CDMX',19.42850000,-99.12770000,'2025-11-15 01:15:28','2025-11-15 01:15:28','[\"martes\", \"jueves\"]',NULL,NULL),(3,'carlos@mail.com','password123','Carlos López','555-0103','donor','Boulevard Norte 789, CDMX',19.44200000,-99.13890000,'2025-11-15 01:15:28','2025-11-15 01:15:28','[\"lunes\", \"martes\", \"miércoles\", \"jueves\", \"viernes\"]',NULL,NULL),(4,'ana@mail.com','password123','Ana Martínez','555-0104','organization','Calle Sur 321, CDMX',19.41950000,-99.14250000,'2025-11-15 01:15:28','2025-11-15 01:15:28','[\"sábado\", \"domingo\"]',NULL,NULL),(5,'pedro@mail.com','password123','Pedro Sánchez','555-0105','donor','Avenida Este 654, CDMX',19.43670000,-99.12080000,'2025-11-15 01:15:28','2025-11-15 01:15:28','[\"lunes\", \"miércoles\", \"viernes\"]',NULL,NULL);
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2025-11-14 21:48:14
