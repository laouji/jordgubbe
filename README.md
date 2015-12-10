# jordgubbe

jordgubbe is an incoming webhook that posts app reviews from the iTunes Store feed to Slack

![sample](https://i.gyazo.com/3e0b33e694eda96be816add2bae6af50.png)

### Installation

1. `go get -d github.com/laouji/jordgubbe`

2. Create a SQLite database using `./middleware/schema.sql`

3. Add your webhook url, iTunes app id and the path to your database to config/config.yml

```yaml
bot_name: "jordgubbe"
icon_emoji: ":strawberry:"
message_text: "Here are the latest reviews for your app:"
web_hook_uri: "https://hooks.slack.com/services/<REST OF THE WEBHOOK URL>"
itunes_app_id: "<YOUR ITUNES STORE APP ID>"
db_path: "/tmp/jordgubbe.db"
max_attachment_count: 5
```
4. `go install`

5. Run periodically using cron or some other job scheduler
