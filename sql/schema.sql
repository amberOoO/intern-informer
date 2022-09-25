CREATE TABLE IF NOT EXISTS job_info (
    id INTEGER PRIMARY KEY,

    created_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    job_name TEXT NOT NULL,
    company TEXT NOT NULL,
    job_description TEXT,
    job_type TEXT,
    job_category TEXT,
    job_location TEXT,

    is_checked Boolean NOT NULL DEFAULT false, -- 判断是否原来就存在的jobInfo, false: 未check, true: 已check。在获取checked所有没checked的info时，顺便置为1
    is_exists Boolean NOT NULL DEFAULT true   -- 判断数据是否已经从官网移除，每次获取新爬虫数据前，全部置为false，插入一条则置true一条。最后为0的即官网删除的。
);
CREATE UNIQUE  INDEX job_info_idx
on job_info (job_name, company);