# archeage-discord-bot

ArcheAge Discord Bot

## Description

이 봇은 디스코드에서 아키에이지와 관련된 데이터를 불러오기 위해 만들어졌습니다. [archeage-go](https://github.com/geeksbaek/archeage-go) 패키지를 사용하며, 빌드를 거치면 로컬에서 구동할 수 있습니다.

## Usage

빌드된 프로그램은 -t 옵션으로 디스코드 봇 토큰을 전달받아 실행됩니다.

`./archeage-discord-bot -t [TOKEN]`

봇이 구동되었다면 아래 URL에 Client ID를 넣어 브라우저에서 접속하고, 해당 페이지에서 봇을 추가할 서버를 선택합니다.

`https://discordapp.com/api/oauth2/authorize?client_id=[CLIENT_ID]&scope=bot&permissions=0`

봇이 추가된 서버에서는 다음 두 가지의 명령어가 지원됩니다.

- 경매장 가격 조회

`?경매장 "검색어" [*[숫자]]`

경매장에서 해당 검색어를 포함하는 아이템의 가격을 가져옵니다.

- 캐릭터 검색

`?캐릭터 "검색어" [@[서버]]`

해당 검색어를 포함하는 캐릭터의 정보를 가져옵니다.