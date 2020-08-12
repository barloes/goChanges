I have stopped using lambda and apigateway since i have found googel stmp api to be much easier to use(check email folder)

create env file with the following variable:

AWS_SECRET_ACCESS_KEY= ??
AWS_SECRET_KEY= ??
REGION= ??
email_user = ??
email_pass = ??

run "go run screenshot.go"

this program aims to track the website through taking screenshots of the website and checking for changes
once there is changes, it will check against the dynamodb db for the email associated to the url and proceed to send email to them to notify them about the changes.

TODO:
i am going to add a frontend and backend to simplify the process of adding url,email to the dynamodb database.