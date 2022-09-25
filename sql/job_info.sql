-- name: InsertJobInfo :exec
INSERT INTO job_info (
    job_name, company, job_description, job_type, job_category, job_location
) VALUES (
    ?, ?, ?, ?, ?, ?
)ON CONFLICT(job_name, company) DO UPDATE SET is_exists = true;

-- name: GetJobInfosByCompany :many
SELECT * FROM job_info
WHERE company = ? LIMIT 1;

-- name: GetJobInfoByJobNameAndCompany :one
SELECT * FROM job_info
WHERE job_name = ? AND company = ? LIMIT 1;

-- name: UpdateIsChecked :exec
UPDATE job_info 
SET is_checked = ?
WHERE id = ?;

-- name: ResetIsExistsByCompany :exec
UPDATE job_info
SET is_exists = 0
WHERE company = ?;

-- name: GetJobInfosByIsExistsAndCompany :many
SELECT * FROM job_info
WHERE is_exists = ? AND company = ?;

-- name: GetJobInfosByIsCheckedAndCompany :many
SELECT * FROM job_info
WHERE is_checked = ? AND company = ?;