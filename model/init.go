package model

import (
	"context"
	"fmt"

	"github.com/dyhes/techtrainingcamp-CourseSelectionSystem/models"
	"github.com/go-redis/redis/v8"
)

var Rdb *redis.Client

var Ctx = context.Background()

func Setup() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "Ndhtbytegroup6..",
		DB:       0,
	})

	Rdb.FlushDB(Ctx) // 清空redis

	var students []Student // 为每一位同学开一个course集合 key = studentcourse%d，并进行课程同步。
	err := models.Db.Find(&students).Error
	fmt.Println(err)
	for _, student := range students {
		Rdb.SAdd(Ctx, fmt.Sprintf("studentcourse%d", student.Studentid), "") // 为每位同学默认加入一个空的课程。
	}
	var studentcourse []StudentCourse
	models.Db.Find(&studentcourse)
	for _, it := range studentcourse {
		Rdb.SAdd(Ctx, fmt.Sprintf("studentcourse%d", it.Studentid), fmt.Sprintf("%d", it.Courseid))
	}

	// 记录每个课程的课程名，任课老师，以及余量
	var courses []Course
	models.Db.Find(&courses)
	for _, course := range courses {
		Rdb.Set(Ctx, fmt.Sprintf("coursename%d", course.Courseid), course.Coursename, redis.KeepTTL)
		Rdb.Set(Ctx, fmt.Sprintf("courseteacher%d", course.Courseid), course.Teacherid, redis.KeepTTL)
		Rdb.Set(Ctx, fmt.Sprintf("course%d", course.Courseid), course.Remaincap, redis.KeepTTL)
	}
	go BookCourseHandler() // 开启一个协程来处理任务
}

/*
create table member (
    user_id bigint primary key auto_increment,
    nickname varchar(30) not null ,
    username varchar(30) not null unique,
    password varchar(30) not null ,
    user_type tinyint not null ,
    is_active tinyint(1)
)auto_increment = 1000000 default charset utf8;

create table student (
    student_id bigint primary key,
    student_name varchar(30) not null unique,
    foreign key (student_id) references member(user_id)
)default charset utf8;

create table teacher (
    teacher_id bigint primary key,
    teacher_name varchar(30) not null unique,
    foreign key (teacher_id) references member(user_id)
)default charset utf8;

create table course (
    course_id bigint primary key auto_increment,
    course_name varchar(30) not null unique,
    cap int unsigned not null,
    remain_cap int unsigned not null,
    teacher_id bigint
)auto_increment = 1000000 default charset utf8;

create table student_course (
    student_id bigint,
    student_name varchar(30) not null,
    course_id bigint,
    course_name varchar(30) not null,
    primary key (student_id, course_id)
)default charset utf8;

# 默认把管理员插入表
insert into member(nickname, username, password, user_type, is_active)
values("JudgeAdmin", "JudgeAdmin", "JudgePassword2022", 1, true);
*/
