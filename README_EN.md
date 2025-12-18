<div align="center">
  <img src="resources/assets/logo.png" alt="Gowlive Logo" width="80"/>
  <h1>Gowlive</h1>
  <p><a href="./CHANGELOG.md">CHANGELOG</a> | <a href="./README.md">中文文档</a></p>
  <p>
    <strong>A live streaming recording platform built with GoFrame and Vue3</strong>
  </p>
  <p>
  <a>
    <img alt="MIT License"
      src="https://img.shields.io/github/license/shichen437/Gowlive">
  </a>
  <a>
    <img alt="Docker Image Version"
      src="https://img.shields.io/docker/v/shichen437/gowlive?labelColor=%20%23FDB062&color=%20%23f79009">
  </a>
  <a href="https://hub.docker.com/u/shichen437" target="_blank">
    <img alt="Docker Pulls"
      src="https://img.shields.io/docker/pulls/shichen437/gowlive?labelColor=%20%23528bff&color=%20%23155EEF">
  </a>
  <a>
    <img alt="Docker Image Size"
      src="https://img.shields.io/docker/image-size/shichen437/gowlive">
  </a>
</p>
</div>

## ✨ Introduction

Gowlive is a live streaming recording platform built with [GoFrame](https://goframe.org/) and [Vue3](https://vuejs.org/), supporting one-click Docker deployment. You can use it to record your favorite live content for easy review.

_This project is for learning and technical exchange purposes only, not for commercial use. It does not involve any private information (including storage, uploading, crawling, etc.)._

## 🚀 Features

- 🔴 **Live Recording**: Supports real-time/timed/smart recording tasks.
- 🕒 **Live History**: Records live history during recording.
- 📽️ **Video Editing**: Supports quick editing of recorded videos.
- 🫥 **Anchor Data**: Records and automatically updates basic anchor information, displays statistical charts.
- 🍪 **Cookie Management**: Conveniently manage login Cookies from different platforms.
- 📁 **File Management**: Manage recorded video files.
- ☑️ **File Check**: Checks if video files are corrupted.
- 📜 **Log Center**: Records system operation logs for easy troubleshooting.
- 📢 **Push Channels**: Supports sending messages through various channels such as Gotify, email.
- 🔔 **Notification Center**: Centralized viewing and management of all notifications.
- 💾 **Storage Alert**: Sends alerts when storage space reaches a threshold.
- 🌙 **Dark Mode**: Supports dark mode.

## 📺 Supported Platforms

- 抖音
- 哔哩哔哩
- YY直播
- Bigo Live

## 📸 Preview

<div align="center">
  <img src="resources/assets/screenshots/en/login.jpg" alt="Login Page" width="92%">
</div>

<br>

<div align="center">
  <img src="resources/assets/screenshots/en/home.jpg" alt="Home Page" width="45%">&nbsp;&nbsp;
  <img src="resources/assets/screenshots/en/settings.jpg" alt="Settings" width="45%">
</div>
<div align="center">
  <img src="resources/assets/screenshots/en/room-card.jpg" alt="Room Card" width="45%">&nbsp;&nbsp;
  <img src="resources/assets/screenshots/en/room-green.jpg" alt="Theme" width="45%">
</div>

## 📦 Deployment

> Initial account password: `admin` / `gowlive`

### Docker Image Deployment

1.  Pull the image:

    ```bash
    docker pull shichen437/gowlive:latest
    ```

2.  Run the container:

    ```bash
    docker run -d \
      -p 12580:12580 \
      -v /data/gowlive:/gowlive/resources/data \
      -e PROJECT_SM4KEY=abcdefghijklmnopqrstuvwxyz123456 \
      -e PROJECT_LANG=en \
      -e TZ=Asia/Shanghai \
      --name gowlive \
      --restart=always \
      shichen437/gowlive:latest
    ```

    > **Note:** `-v /data/gowlive:/gowlive/resources/data` is used for persistent storage of recorded files and application data. Please replace `/data/gowlive` with your actual local path.

3.  Visit `http://<YOUR_IP>:12580` to view the application.

### Docker Compose Deployment

1.  Download and rename the `docker-compose.yaml.example` file:
    ```bash
    wget https://raw.githubusercontent.com/shichen437/Gowlive/main/docker-compose.yaml.example -O docker-compose.yaml
    ```
2.  Modify the `docker-compose.yaml` file as needed.

    **Environment Variables:**

| Variable Name    | Description                        | Default Value                      | Required |
| :--------------- | :--------------------------------- | :--------------------------------- | :------- |
| `PROJECT_SM4KEY` | SM4 Encryption Key (32-bit string) | `abcdefghijklmnopqrstuvwxyz123456` | Yes      |
| `TZ`             | Timezone                           | `Asia/Shanghai`                    | No       |
| `PROJECT_LANG`   | Language                           | `zh-CN` (`en`, `zh-TW`)            | No       |

3.  Start the service:
    ```bash
    docker-compose up -d
    ```
4.  Visit `http://<YOUR_IP>:12580` to view the application.

## 🤝 Acknowledgements

Sincere thanks to every contributor for their support and dedication to Gowlive. Your contributions are greatly appreciated.

(<a href="./.github/docs/CONTRIBUTING_EN.md">Contributing Guide</a>)

<a href="https://github.com/shichen437/Gowlive/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=shichen437/Gowlive" />
</a>

## 📄 License

This project is open source under the [MIT License](https://github.com/shichen437/Gowlive/blob/main/LICENSE).
