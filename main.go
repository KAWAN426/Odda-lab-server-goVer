package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"

	middleware "oddalab_server/Tools/middleware"

	"github.com/lib/pq"
)

type Objects struct {
	Id   string
	Name string
	Img  string
}

type Lab struct {
	Id            string                 `gorm:"column:id"`
	Title         string                 `gorm:"column:title"`
	MakerID       string                 `gorm:"column:makerId"`
	Objects       []Objects              `gorm:"column:objects;type:jsonb"`
	BackgroundImg string                 `gorm:"column:backgroundImg"`
	Combinate     [][]string             `gorm:"column:combinate;type:text[][]"`
	StartObj      []string               `gorm:"column:startObj;type:text[]"`
	EndObj        []string               `gorm:"column:endObj;type:text[]"`
	LikedUser     []string               `gorm:"column:likedUser;type:text[]"`
	FindObj       map[string]interface{} `gorm:"column:findObj;type:jsonb"`
	CreatedAt     string                 `gorm:"column:createdAt"`
	UpdatedAt     string                 `gorm:"column:updatedAt"`
}

func main() {
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	middleware.Setting(app)

	connStr := "user=postgres password=0426 dbname=main sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	app.Get("/lab/popular", func(c *fiber.Ctx) error {
		rows, err := db.Query("SELECT id,title,makerId,objects,backgroundImg,combinate,startObj,endObj,likedUser,findObj,createdAt,updatedAt FROM lab ORDER BY ARRAY_LENGTH(likedUser, 1) DESC")
		var labs []Lab
		for rows.Next() {
			var lab Lab
			var objects_data []uint8
			var findObj_data []uint8
			var combinate_data []uint8
			err := rows.Scan(
				&lab.Id,
				&lab.Title,
				&lab.MakerID,
				&objects_data,
				&lab.BackgroundImg,
				&combinate_data,
				pq.Array(&lab.StartObj),
				pq.Array(&lab.EndObj),
				pq.Array(&lab.LikedUser),
				&findObj_data,
				&lab.CreatedAt,
				&lab.UpdatedAt,
			)
			if err != nil {
				log.Fatal(err)
			}
			json.Unmarshal(objects_data, &lab.Objects)
			json.Unmarshal(findObj_data, &lab.FindObj)
			lab.Combinate = convertUint8To2DArray(combinate_data)
			labs = append(labs, lab)
		}
		if err != nil {
			return c.Status(400).JSON(err.Error())
		}
		return c.Status(200).JSON(labs)
	})
	app.Get("/lab/newest", func(c *fiber.Ctx) error {
		rows, err := db.Query("SELECT id,title,makerId,objects,backgroundImg,combinate,startObj,endObj,likedUser,findObj,createdAt,updatedAt FROM lab ORDER BY createdAt DESC")
		var labs []Lab
		for rows.Next() {
			var lab Lab
			var objects_data []uint8
			var findObj_data []uint8
			var combinate_data []uint8
			err := rows.Scan(
				&lab.Id,
				&lab.Title,
				&lab.MakerID,
				&objects_data,
				&lab.BackgroundImg,
				&combinate_data,
				pq.Array(&lab.StartObj),
				pq.Array(&lab.EndObj),
				pq.Array(&lab.LikedUser),
				&findObj_data,
				&lab.CreatedAt,
				&lab.UpdatedAt,
			)
			if err != nil {
				log.Fatal(err)
			}
			json.Unmarshal(objects_data, &lab.Objects)
			json.Unmarshal(findObj_data, &lab.FindObj)
			lab.Combinate = convertUint8To2DArray(combinate_data)
			labs = append(labs, lab)
		}
		if err != nil {
			return c.Status(400).JSON(err.Error())
		}
		return c.Status(200).JSON(labs)
	})

	app.Listen(":5000")
}

func convertUint8To2DArray(data []uint8) [][]string {
	re := regexp.MustCompile(`[{}]+`)
	str := re.ReplaceAllString(string(data), "")
	strSlice := strings.Split(str, "\n")
	stringSliceSlice := make([][]string, len(strSlice))

	for i, s := range strSlice {
		items := strings.Split(strings.TrimSpace(s), ",")
		stringSliceSlice[i] = make([]string, len(items))

		for j, item := range items {
			stringSliceSlice[i][j] = strings.TrimSpace(item)
		}
	}

	return stringSliceSlice
}
