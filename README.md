# SE69 Lab04 Starter: Lab Equipment Loan API

Go 1.24 ขึ้นไป ใช้ Gin (web framework), GORM (ORM) และ SQLite (ไม่ต้องติดตั้งโปรแกรมฐานข้อมูลเพิ่ม)

## วิธีรัน

```bash
cp .env.example .env # Windows: copy .env.example .env  (แล้วแก้ JWT_SECRET)
go mod tidy          # ดาวน์โหลด library ครั้งแรก (ต้องต่ออินเทอร์เน็ต)
go run ./cmd/api     # API อยู่ที่ http://localhost:8080
```

ทุกครั้งที่รัน API ข้อมูลในไฟล์ `lab04.db` จะถูกล้างและเริ่มต้นใหม่

ค่าตั้งค่าอยู่ในไฟล์ `.env` (โปรแกรมอ่านให้เองตอนเริ่ม) หลัง clone ให้คัดลอก `.env.example` เป็น `.env` แล้วเปลี่ยน `JWT_SECRET` เป็นข้อความสุ่มยาว ๆ ของตัวเอง ไฟล์ `.env` เก็บ secret จึงห้าม commit หรือแชร์ (`.gitignore` กันไว้แล้ว)

| ตัวแปร       | ค่าใน `.env`   | ความหมาย                     |
|--------------|--------------|------------------------------|
| `PORT`       | `8080`       | port ของ API                 |
| `DB_PATH`    | `lab04.db`   | ไฟล์ SQLite                   |
| `TOKEN_TTL`  | `15m`        | อายุของ token เช่น `30s`       |
| `JWT_SECRET` | (ตั้งเอง)      | key สำหรับเซ็น token (HS256) ถ้าไม่ได้ตั้งค่า server จะไม่เริ่มทำงาน |

## Library ที่ใช้

| Library | ใช้ทำอะไร |
|---------|----------|
| `github.com/gin-gonic/gin` | รับ request, จับคู่ route, middleware, อ่าน JSON/form, ตั้ง cookie |
| `gorm.io/gorm` | ORM: สร้างตารางจาก struct และ query แบบ parameterized (กัน SQL injection) |
| `github.com/ncruces/go-sqlite3/gormlite` | SQLite driver ของ GORM เขียนด้วย Go ล้วน ไม่ต้องใช้ cgo หรือ gcc |
| `github.com/golang-jwt/jwt/v5` | สร้างและตรวจ JWT (ตรวจ alg, signature, iss, aud, exp) |
| `github.com/joho/godotenv` | อ่านค่าตั้งค่าจากไฟล์ `.env` |
| `golang.org/x/crypto/bcrypt` | hash รหัสผ่านแบบมี salt และช้าโดยตั้งใจ |
| `crypto/rand` (standard library) | สุ่ม session ID ที่เดาไม่ได้ |

## ผู้ใช้ทดลอง (รหัสผ่าน `pass1234` ทุกคน)

| Username   | Role    |
|------------|---------|
| `student1` | student |
| `student2` | student (ยืมอุปกรณ์ id 1 อยู่ เป็นรายการยืม id 1) |
| `staff1`   | staff   |
| `teacher1` | teacher |

## Endpoints

| Method และ path               | ใครเรียกได้      |
|-------------------------------|----------------|
| `GET /health`                 | ทุกคน          |
| `POST /sessions` (HTTP Basic) | ทุกคน          |
| `GET /equipment?category=`    | ผู้ที่มี token    |
| `GET /equipment/{id}`         | ผู้ที่มี token    |
| `POST /equipment`             | staff          |
| `PUT /equipment/{id}`         | staff          |
| `GET /deleteEquipment/{id}`   | staff          |
| `POST /loans`                 | student        |
| `GET /loans`                  | ผู้ที่มี token    |
| `GET /loans/{id}`             | ผู้ที่มี token    |
| `GET /reports/usage`          | teacher        |

### Web app (ใช้ session cookie แทน Bearer token)

| Method และ path                                   | ใครเรียกได้        |
|---------------------------------------------------|------------------|
| `POST /web/sessions` (form: `username`, `password`) | ทุกคน (login)     |
| `DELETE /web/sessions`                            | ทุกคน (logout)    |
| `GET /web/loans`                                  | ผู้ที่มี session cookie |
| `POST /web/loans` (form: `equipmentId`)           | ผู้ที่มี session cookie |

รายละเอียด request และ response ดูได้ที่ `openapi.yaml`

## โครงสร้างโปรเจกต์

```
cmd/api/main.go          จุดเริ่มโปรแกรม: โหลด config, เปิดฐานข้อมูล, เริ่ม server
internal/config/         อ่านค่าตั้งค่าจาก .env และ environment variable
internal/models/         struct ของตาราง (GORM) และการสร้างตาราง + ข้อมูลตัวอย่าง (database.go)
internal/dto/            รูปแบบข้อมูลที่รับเข้า (request พร้อม binding) และส่งออก (response)
internal/controllers/    รับ request, query ด้วย GORM, ส่ง response
internal/middleware/     ตรวจ Bearer token และ role (auth.go), ตรวจ session cookie (session.go), log request (logger.go)
internal/routes/         จับคู่ method + path กับ controller และกำหนด role ที่เรียกได้
internal/utils/          ฟังก์ชันช่วย: JWT, hash รหัสผ่าน, สุ่ม session id, ส่ง error response
```
