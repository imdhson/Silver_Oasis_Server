[대구대학교 공학제 Silver_Oasis의 서버사이드] <br>
대구대학교 공학제 Silver_Oasis의 서버사이드에 쓰인 소스코드입니다

설치 시 해야할 것 : <br>
pip install bardapi <br>
pip install load_dotenv <br>
go언어 환경에서 godotenv 필요 <br>
systemctl 사용시 위의 것을 sudo 로 설치해야함 <br>


<br> 본 디렉터리의 .env 파일에<br>
 MongoDB서버, smpt서버 Password, BardAPI key 를 넣어야 동작함.
 <br>sslforfree 폴더에 (본 디렉터리 하위에 없으면 생성해야함)
 <br> https://blog.naver.com/imdhson/223142347090 (손동휘 작성 블로그 글)<br>
 참조하여 ca_bundle.crt, certificate.crt(combined.crt), private.key 를 넣어야 함.


<br> mongoDB 관련<br>
dump 를 mongodb의 mongorestore 명령어를 이용하여 사용하고자 하는 MongoDB서버에 restore하면 시연과 같은 효과를 볼 수 있음. 단, 데이터셋의 출처를 명시해야함.(보건복지부) 
