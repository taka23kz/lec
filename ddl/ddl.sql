-- ユーザテーブル
-- 講義を作成する親ユーザ
create table t_user 
(
	id varchar(30),
	user_name varchar(200) not null,
	mail_address varchar(300) not null,
	user_status varchar(2) not null, -- 00:仮登録状態,10:登録状態,20:停止状態
	limit_flag boolean,  -- true:有料機能制限あり, false:有料機能制限なし
	passwd text not null,
	delete_flag boolean, -- true:無効データ, false:有効データ
	Created timestamp DEFAULT CURRENT_TIMESTAMP,
	Updated timestamp DEFAULT CURRENT_TIMESTAMP,
	primary key (user_id)
);

-- 受講者テーブル
-- 講義を受けるユーザ用のテーブル
create table t_student
(
	id varchar(30),
	owner_user_id varchar(30) references t_user(id),
	family_name varchar(100) not null,
	first_name varchar(100) not null,
	mail_address varchar(300),
	passwd text not null,
	delete_flag boolean, -- true:無効データ, false:有効データ
	Created timestamp DEFAULT CURRENT_TIMESTAMP,
	Updated timestamp DEFAULT CURRENT_TIMESTAMP,
	primary key (id, owner_user_id)
);

-- グループテーブル
create table t_group
(
	id SERIAL not null,
	group_name varchar(300) not null,
	owner_user_id varchar(30) references t_user(id),
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



