USE [master]
GO
/****** Object:  Database [lbe]    Script Date: 16/4/2025 11:18:30 PM ******/
CREATE DATABASE [lbe]
 CONTAINMENT = NONE
 ON  PRIMARY 
( NAME = N'lbe', FILENAME = N'/var/opt/mssql/data/lbe.mdf' , SIZE = 8192KB , MAXSIZE = UNLIMITED, FILEGROWTH = 65536KB )
 LOG ON 
( NAME = N'lbe_log', FILENAME = N'/var/opt/mssql/data/lbe_log.ldf' , SIZE = 8192KB , MAXSIZE = 2048GB , FILEGROWTH = 65536KB )
 WITH CATALOG_COLLATION = DATABASE_DEFAULT
GO
ALTER DATABASE [lbe] SET COMPATIBILITY_LEVEL = 150
GO
IF (1 = FULLTEXTSERVICEPROPERTY('IsFullTextInstalled'))
begin
EXEC [lbe].[dbo].[sp_fulltext_database] @action = 'enable'
end
GO
ALTER DATABASE [lbe] SET ANSI_NULL_DEFAULT OFF 
GO
ALTER DATABASE [lbe] SET ANSI_NULLS OFF 
GO
ALTER DATABASE [lbe] SET ANSI_PADDING OFF 
GO
ALTER DATABASE [lbe] SET ANSI_WARNINGS OFF 
GO
ALTER DATABASE [lbe] SET ARITHABORT OFF 
GO
ALTER DATABASE [lbe] SET AUTO_CLOSE OFF 
GO
ALTER DATABASE [lbe] SET AUTO_SHRINK OFF 
GO
ALTER DATABASE [lbe] SET AUTO_UPDATE_STATISTICS ON 
GO
ALTER DATABASE [lbe] SET CURSOR_CLOSE_ON_COMMIT OFF 
GO
ALTER DATABASE [lbe] SET CURSOR_DEFAULT  GLOBAL 
GO
ALTER DATABASE [lbe] SET CONCAT_NULL_YIELDS_NULL OFF 
GO
ALTER DATABASE [lbe] SET NUMERIC_ROUNDABORT OFF 
GO
ALTER DATABASE [lbe] SET QUOTED_IDENTIFIER OFF 
GO
ALTER DATABASE [lbe] SET RECURSIVE_TRIGGERS OFF 
GO
ALTER DATABASE [lbe] SET  DISABLE_BROKER 
GO
ALTER DATABASE [lbe] SET AUTO_UPDATE_STATISTICS_ASYNC OFF 
GO
ALTER DATABASE [lbe] SET DATE_CORRELATION_OPTIMIZATION OFF 
GO
ALTER DATABASE [lbe] SET TRUSTWORTHY OFF 
GO
ALTER DATABASE [lbe] SET ALLOW_SNAPSHOT_ISOLATION OFF 
GO
ALTER DATABASE [lbe] SET PARAMETERIZATION SIMPLE 
GO
ALTER DATABASE [lbe] SET READ_COMMITTED_SNAPSHOT OFF 
GO
ALTER DATABASE [lbe] SET HONOR_BROKER_PRIORITY OFF 
GO
ALTER DATABASE [lbe] SET RECOVERY FULL 
GO
ALTER DATABASE [lbe] SET  MULTI_USER 
GO
ALTER DATABASE [lbe] SET PAGE_VERIFY CHECKSUM  
GO
ALTER DATABASE [lbe] SET DB_CHAINING OFF 
GO
ALTER DATABASE [lbe] SET FILESTREAM( NON_TRANSACTED_ACCESS = OFF ) 
GO
ALTER DATABASE [lbe] SET TARGET_RECOVERY_TIME = 60 SECONDS 
GO
ALTER DATABASE [lbe] SET DELAYED_DURABILITY = DISABLED 
GO
ALTER DATABASE [lbe] SET ACCELERATED_DATABASE_RECOVERY = OFF  
GO
EXEC sys.sp_db_vardecimal_storage_format N'lbe', N'ON'
GO
ALTER DATABASE [lbe] SET QUERY_STORE = OFF
GO
USE [lbe]
GO
/****** Object:  Table [dbo].[sys_channel]    Script Date: 16/4/2025 11:18:30 PM ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[sys_channel](
	[id] [bigint] IDENTITY(1,1) NOT NULL,
	[app_id] [varchar](100) NOT NULL,
	[app_key] [varchar](100) NOT NULL,
	[status] [char](2) NOT NULL,
	[sig_method] [varchar](100) NOT NULL,
	[create_time] [datetime] NULL,
	[update_time] [datetime] NOT NULL,
PRIMARY KEY CLUSTERED 
(
	[id] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON, OPTIMIZE_FOR_SEQUENTIAL_KEY = OFF) ON [PRIMARY]
) ON [PRIMARY]
GO
/****** Object:  Table [dbo].[users]    Script Date: 16/4/2025 11:18:30 PM ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[users](
	[id] [bigint] IDENTITY(1,1) NOT NULL,
	[external_id] [varchar](50) NOT NULL,
	[opted_in] [bit] NOT NULL,
	[external_id_type] [varchar](50) NULL,
	[email] [varchar](50) NULL,
	[dob] [date] NULL,
	[country] [varchar](50) NULL,
	[first_name] [varchar](255) NULL,
	[last_name] [varchar](255) NULL,
	[burn_pin] [int] NULL,
	[created_at] [datetime] NULL,
	[updated_at] [datetime] NULL,
PRIMARY KEY CLUSTERED 
(
	[id] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON, OPTIMIZE_FOR_SEQUENTIAL_KEY = OFF) ON [PRIMARY]
) ON [PRIMARY]
GO
/****** Object:  Table [dbo].[users_phone_numbers]    Script Date: 16/4/2025 11:18:30 PM ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[users_phone_numbers](
	[id] [bigint] IDENTITY(1,1) NOT NULL,
	[user_id] [bigint] NOT NULL,
	[phone_number] [varchar](20) NULL,
	[phone_type] [varchar](20) NULL,
	[preference_flags] [varchar](50) NULL,
	[created_at] [datetime] NULL,
	[updated_at] [datetime] NULL,
PRIMARY KEY CLUSTERED 
(
	[id] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON, OPTIMIZE_FOR_SEQUENTIAL_KEY = OFF) ON [PRIMARY]
) ON [PRIMARY]
GO
ALTER TABLE [dbo].[sys_channel] ADD  DEFAULT ('10') FOR [status]
GO
ALTER TABLE [dbo].[sys_channel] ADD  DEFAULT ('SHA256') FOR [sig_method]
GO
ALTER TABLE [dbo].[users] ADD  DEFAULT ((0)) FOR [opted_in]
GO
ALTER TABLE [dbo].[users] ADD  DEFAULT (getdate()) FOR [created_at]
GO
ALTER TABLE [dbo].[users] ADD  DEFAULT (getdate()) FOR [updated_at]
GO
ALTER TABLE [dbo].[users_phone_numbers] ADD  DEFAULT (getdate()) FOR [created_at]
GO
ALTER TABLE [dbo].[users_phone_numbers] ADD  DEFAULT (getdate()) FOR [updated_at]
GO
ALTER TABLE [dbo].[users_phone_numbers]  WITH CHECK ADD  CONSTRAINT [fk_user_id] FOREIGN KEY([user_id])
REFERENCES [dbo].[users] ([id])
ON DELETE CASCADE
GO
ALTER TABLE [dbo].[users_phone_numbers] CHECK CONSTRAINT [fk_user_id]
GO
USE [master]
GO
ALTER DATABASE [lbe] SET  READ_WRITE 
GO
