-- ユーザテーブル
create table t_user 
(
	id SERIAL not null, 
	user_id varchar(30),
	user_name varchar(200) not null,
	mail_address varchar(300) not null,
	user_type varchar(2) not null, 	-- 00:管理者,10:受講者,99:Guest
	user_status varchar(2) not null, -- 00:仮登録状態,10:登録状態,20:停止状態
	limit_flag boolean,  -- true:有料機能制限あり, false:有料機能制限なし
	passwd text not null,
	primary key (id)
);

-- 組織テーブル
create table t_org
(
	id SERIAL not null, 
	org_name varchar(200) unique not null,
	parent_org_id int references t_org(id),
	primary key (id)
)

-- ユーザと組織を紐づけるためのテーブル
create table t_user_org_rel
(
	user_id int references t_user(id),
	org_id int references t_org(id)
)

-- 講義のオーナーが所属するグループ(組織)のテーブル
create table t_owner_group
(
	owner_group_id SERIAL not null,
	owner_id int not null references t_owner(owner_id),
	owner_group_name varchar(200) not null,
	parent_owner_group_id int references t_owner_group(owner_group_id),
	primary key (owner_group_id)
);

-- 講義テーブル
create table t_lesson
(
	lesson_id SERIAL not null,
	lesson_name varchar(300),
	sub_number int not null,		-- 枝番
	content_type varchar(2),		-- 00:movie, 01:slide, 02:text, 03:quiz
	parent_lesson_id int references t_lesson(lesson_id),
	owner_group_id int not null references t_owner_group(owner_group_id),
	primary key (lesson_id)
);

-- 設問テーブル（多肢選択系、回答入力系の問題を管理するためのテーブル)
create table t_question 
(
	question_id SERIAL not null,
	question text not null,					-- 設問
	answer_type varchar(2) not null ,		-- 00:text, 01:radio, 02:checkbox
	choice_num integer not null,			-- answer_typeがradio,checkboxの場合に選択肢をいくつ表示するか。
	owner_group_id int not null references t_owner_group(owner_group_id),
	lesson_id int not null references t_lesson(lesson_id),
	primary key (question_id)
);

-- 選択肢テーブル（t_questionに紐づく選択肢を管理するテーブル)
create table t_choice
(
	choice_id SERIAL not null,
	question_id int not null references t_question(question_id),
	choice_label text not null,										-- 選択肢
	correct boolean not null,										-- true:正解 false:不正解
	primary key (choice_id)
);







