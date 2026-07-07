-- phpMyAdmin SQL Dump
-- version 5.2.1
-- https://www.phpmyadmin.net/
--
-- Host: 127.0.0.1
-- Generation Time: Jul 07, 2026 at 07:21 PM
-- Server version: 10.4.32-MariaDB
-- PHP Version: 8.0.30

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `usr_d`
--

-- --------------------------------------------------------

--
-- Table structure for table `__story_tale__`
--

CREATE TABLE `__story_tale__` (
  `__SID__` int(11) NOT NULL,
  `__USID__` int(11) NOT NULL,
  `__TITILE__` varchar(100) NOT NULL,
  `__DESCRIPTION__` text NOT NULL,
  `__THUMBNAIL__` varchar(255) NOT NULL,
  `__GENRE__` varchar(50) NOT NULL,
  `__VISIBILITY__` enum('PUBLIC','PRIVATE') NOT NULL,
  `__SCENE_COUNT__` int(11) NOT NULL,
  `__CREATED__` timestamp NOT NULL DEFAULT current_timestamp(),
  `__UPDATED__` timestamp NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table `__usr_tale__`
--

CREATE TABLE `__usr_tale__` (
  `__USID__` int(11) NOT NULL,
  `__USER_NAME__` varchar(50) NOT NULL,
  `__EMAIL__` varchar(100) NOT NULL,
  `__PASSWORD__` varchar(255) NOT NULL,
  `__P_PIC__` varchar(255) NOT NULL,
  `__ORIGIN__` timestamp NOT NULL DEFAULT current_timestamp(),
  `__LAST_SEEN__` timestamp NOT NULL DEFAULT current_timestamp(),
  `__STATUS__` enum('Active','Blocked') NOT NULL DEFAULT 'Active'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Indexes for dumped tables
--

--
-- Indexes for table `__story_tale__`
--
ALTER TABLE `__story_tale__`
  ADD PRIMARY KEY (`__SID__`);

--
-- Indexes for table `__usr_tale__`
--
ALTER TABLE `__usr_tale__`
  ADD PRIMARY KEY (`__USID__`),
  ADD UNIQUE KEY `__USER_NAME__` (`__USER_NAME__`),
  ADD UNIQUE KEY `__EMAIL__` (`__EMAIL__`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `__story_tale__`
--
ALTER TABLE `__story_tale__`
  MODIFY `__SID__` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `__usr_tale__`
--
ALTER TABLE `__usr_tale__`
  MODIFY `__USID__` int(11) NOT NULL AUTO_INCREMENT;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
