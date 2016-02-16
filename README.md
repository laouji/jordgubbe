# jordgubbe

jordgubbe is an incoming webhook that posts app reviews from the iTunes Store and Google Play to Slack.

To retrieve android reviews you will need to have gsutil installed.

![sample](https://i.gyazo.com/3e0b33e694eda96be816add2bae6af50.png)

### Installation

1. `go get github.com/laouji/jordgubbe`

2. Create a SQLite database  `sqlite3 /tmp/jordgubbe.db < ./middleware/schema.sql`

3. Create a config.yml file similar to the one below, adding your webhook url, iTunes app id and the path to your database, etc

```yaml
bot_name: "jordgubbe"
icon_emoji: ":strawberry:"
message_text: "Here are the latest reviews for your app:"
web_hook_uri: "https://hooks.slack.com/services/<REST OF THE WEBHOOK URL>"
itunes_app_id: "<YOUR ITUNES STORE APP ID>"
db_path: "/tmp/jordgubbe.db"
tmp_dir: "/tmp"
gcs_bucket_id: "<YOUR GOOGLE CLOUD STORAGE BUCKET ID>"
android_package_name: "com.example.appname"
max_attachment_count: 5
```
4. Run periodically using cron or some other job scheduler

```
~/go/bin/jordgubbe -c ~/config.yml -p ios
~/go/bin/jordgubbe -c ~/config.yml -p android

```
