

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for short_link
-- ----------------------------
DROP TABLE IF EXISTS `short_link`;
CREATE TABLE `short_link`  (
  `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `key` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT 'ID 的base62编码串',
  `url` tinytext CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT 'URL地址',
  `create_time` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0),
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
