# Postman Test Suite

## Import
1. Import `Applicant_API.postman_collection.json`.
2. Import `Local.postman_environment.json`.
3. Pilih environment **Sekolah Karir - Local**.
4. Jalankan API dengan `go run .`.
5. Buka collection dan klik **Run**. Jalankan request dalam urutan yang tersedia.

## Test coverage
Collection berisi 23 request dengan automated assertions: health, create, detail, list, search, filter stage/position, pagination, update stage, summary, required field validation, invalid email, duplicate email, invalid stage, invalid query, not found, dan invalid ID.

Request Create membuat email unik memakai timestamp dan menyimpan ID/email sebagai collection variables. Karena itu request lanjutan tidak bergantung pada ID seed yang tetap.

## Newman
```bash
npm install -g newman
newman run Applicant_API.postman_collection.json -e Local.postman_environment.json
```
