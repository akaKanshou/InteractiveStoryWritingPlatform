-- phpMyAdmin SQL Dump
-- version 5.2.1
-- https://www.phpmyadmin.net/
--
-- Host: 127.0.0.1
-- Generation Time: Jul 11, 2026 at 08:06 PM
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
-- Table structure for table `__edge_tale__`
--

CREATE TABLE `__edge_tale__` (
  `__EDGE_ID__` int(11) NOT NULL,
  `__STORY_ID__` int(11) NOT NULL,
  `__FROM_SCENE__` int(11) NOT NULL,
  `__TO_SCENE__` int(11) NOT NULL,
  `__CHOICE_TEXT__` varchar(255) NOT NULL,
  `__CREATED__` timestamp NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table `__scene_tale__`
--

CREATE TABLE `__scene_tale__` (
  `__SCENE_ID__` int(11) NOT NULL,
  `__STORY_ID__` int(11) NOT NULL,
  `__SCENE_NAME__` varchar(100) NOT NULL,
  `__SCENE_TEXT__` longtext DEFAULT NULL,
  `__POS_X__` int(11) DEFAULT 0,
  `__POS_Y__` int(11) DEFAULT 0,
  `__CREATED__` timestamp NOT NULL DEFAULT current_timestamp(),
  `__UPDATED__` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

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

--
-- Dumping data for table `__story_tale__`
--

INSERT INTO `__story_tale__` (`__SID__`, `__USID__`, `__TITILE__`, `__DESCRIPTION__`, `__THUMBNAIL__`, `__GENRE__`, `__VISIBILITY__`, `__SCENE_COUNT__`, `__CREATED__`, `__UPDATED__`) VALUES
(1, 1, 'New Story', 'A story', '290f41b286e0b4a66d1fc83371b850e0.png', 'Fantasy', 'PUBLIC', 0, '2026-07-11 15:49:48', '2026-07-11 15:49:48');

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
  `__STATUS__` enum('Active','Blocked') NOT NULL DEFAULT 'Active',
  `__AUTH_TYPE__` enum('LOCAL','GOOGLE') NOT NULL,
  `__GOOGLE_ID__` varchar(255) NOT NULL,
  `__EMAIL__VERIFIED__` tinyint(1) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `__usr_tale__`
--

INSERT INTO `__usr_tale__` (`__USID__`, `__USER_NAME__`, `__EMAIL__`, `__PASSWORD__`, `__P_PIC__`, `__ORIGIN__`, `__LAST_SEEN__`, `__STATUS__`, `__AUTH_TYPE__`, `__GOOGLE_ID__`, `__EMAIL__VERIFIED__`) VALUES
(1, 'TheCertifiedBerozgar', 'genius22308a@gmail.com', '$2y$10$sUyDCoY4tKoKexSlIE7NRe/TbfcdhFR4T1E9ktS.JixXAJVtQridG', '', '2026-07-11 15:48:56', '2026-07-11 15:56:34', 'Active', 'LOCAL', '', 0);

--
-- Indexes for dumped tables
--

--
-- Indexes for table `__edge_tale__`
--
ALTER TABLE `__edge_tale__`
  ADD PRIMARY KEY (`__EDGE_ID__`),
  ADD KEY `__STORY_ID__` (`__STORY_ID__`),
  ADD KEY `__FROM_SCENE__` (`__FROM_SCENE__`),
  ADD KEY `__TO_SCENE__` (`__TO_SCENE__`);

--
-- Indexes for table `__scene_tale__`
--
ALTER TABLE `__scene_tale__`
  ADD PRIMARY KEY (`__SCENE_ID__`),
  ADD KEY `__STORY_ID__` (`__STORY_ID__`);

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
-- AUTO_INCREMENT for table `__edge_tale__`
--
ALTER TABLE `__edge_tale__`
  MODIFY `__EDGE_ID__` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `__scene_tale__`
--
ALTER TABLE `__scene_tale__`
  MODIFY `__SCENE_ID__` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `__story_tale__`
--
ALTER TABLE `__story_tale__`
  MODIFY `__SID__` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=2;

--
-- AUTO_INCREMENT for table `__usr_tale__`
--
ALTER TABLE `__usr_tale__`
  MODIFY `__USID__` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=2;

--
-- Constraints for dumped tables
--

--
-- Constraints for table `__edge_tale__`
--
ALTER TABLE `__edge_tale__`
  ADD CONSTRAINT `__edge_tale___ibfk_1` FOREIGN KEY (`__STORY_ID__`) REFERENCES `__story_tale__` (`__SID__`) ON DELETE CASCADE,
  ADD CONSTRAINT `__edge_tale___ibfk_2` FOREIGN KEY (`__FROM_SCENE__`) REFERENCES `__scene_tale__` (`__SCENE_ID__`) ON DELETE CASCADE,
  ADD CONSTRAINT `__edge_tale___ibfk_3` FOREIGN KEY (`__TO_SCENE__`) REFERENCES `__scene_tale__` (`__SCENE_ID__`) ON DELETE CASCADE;

--
-- Constraints for table `__scene_tale__`
--
ALTER TABLE `__scene_tale__`
  ADD CONSTRAINT `__scene_tale___ibfk_1` FOREIGN KEY (`__STORY_ID__`) REFERENCES `__story_tale__` (`__SID__`) ON DELETE CASCADE;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
