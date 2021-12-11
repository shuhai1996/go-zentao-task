CREATE TABLE `blog_article` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(45) NOT NULL DEFAULT '' COMMENT '文章标题',
  `author` varchar(100) NOT NULL COMMENT '作者openid',
  `desc` varchar(255) NOT NULL DEFAULT '' COMMENT '文章简述',
  `content` text NOT NULL COMMENT '文章内容',
  `likes` int(11) NOT NULL DEFAULT '0' COMMENT '点赞数',
  `cover_image_url` varchar(255) NOT NULL DEFAULT '' COMMENT '封面图',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '1 正常 2 删除',
  `create_time` datetime NOT NULL,
  `update_time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `author` (`author`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='博客文章'