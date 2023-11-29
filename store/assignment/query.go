package assignment

const GetQuery = "SELECT id, name, points, num_of_attemps, deadline, assignment_created, assignment_updated FROM assignments"

const GetByIDQuery = "SELECT id, user_id, name, points, num_of_attemps, deadline, assignment_created, assignment_updated FROM assignments WHERE id=$1"

const IsAssignmentExistsQuery = "SELECT user_id FROM assignments WHERE id=$1"

const InsertQuery = "INSERT INTO assignments (id, user_id, name, points, num_of_attemps, deadline, assignment_created, assignment_updated) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"

const UpdateQuery = "UPDATE assignments SET name=$1, points=$2, num_of_attemps=$3, deadline=$4, assignment_updated=$5 WHERE id=$6"

const DeleteQuery = "DELETE FROM assignments WHERE id=$1"

const DeleteAssignmentSubmissions = "DELETE FROM submissions WHERE assignment_id=$1"

const GetAssignmentSubmissionQuery = "SELECT num_of_attemps, deadline FROM assignments WHERE id=$1"

// Assignment Submission Queries
const CheckSubmissionsQuery = "SELECT COUNT(*) as submission_count FROM submissions WHERE assignment_id=$1 AND user_id=$2"

const InsertSubmissionQuery = "INSERT INTO submissions(id, user_id, assignment_id, submission_url, submission_created, submission_updated) VALUES ($1, $2, $3, $4, $5, $6)"
