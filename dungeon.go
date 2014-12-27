package main

import (
  "fmt"
  "math/rand"
  "time"
)

// dungeon size
var dungeonWidth = 40
var dungeonHeight = 40

// dungeon rooms
var roomAttempts = 200
var minRoomSize = 5
var maxRoomSize = 10

const (
  WALL = 0
  FLOOR = 1
  EDGE = 2
  DOOR = 3
  TUNNEL = 4
)

func createEmptyDungeon(width int, height int) [][]int {
  fmt.Println("Creating empty dungeon...")
  dungeon := make([][]int, height)
  for i := range dungeon {
    dungeon[i] = make([]int, width)
  }

  return dungeon
}

type Point struct {
  x int
  y int
}

type Room struct {
  width int
  height int
  x int
  y int
  edges []Point
}

type Region struct {
  tiles []Point
  edges []Point
}

func createRooms(dungeon [][]int, minSize, maxSize, attempts int) []Room {
  fmt.Println("Creating rooms...")
  var rooms []Room

  for i := 0; i < attempts; i++ {
    width := rand.Intn(maxSize - minSize) + minSize
    height := rand.Intn(maxSize - minSize) + minSize

    maxX := len(dungeon[0]) - width - 2
    maxY := len(dungeon) - height - 2

    x := rand.Intn(maxX - 3) + 3
    y := rand.Intn(maxY - 3) + 3

    shouldAppend := true
    for r := range rooms {
      if(x + width < rooms[r].x || // to the left
      x > rooms[r].x + rooms[r].width || // to the right
      y + height < rooms[r].y || // fully above
      y > rooms[r].y + rooms[r].height) { // fully below
        // do nothing
      } else {
        shouldAppend = false
        break
      }
    }

    if(shouldAppend) {
      rooms = append(rooms, Room{ width: width, height: height, x: x, y: y })
    }
  }

  for r := range rooms {
    for i := rooms[r].x; i < rooms[r].x + rooms[r].width; i++ {
      for j := rooms[r].y; j < rooms[r].y + rooms[r].height; j++ {
        dungeon[j][i] = FLOOR
      }
    }
  }

  return rooms
}

func createMaze(dungeon [][]int, rooms []Room) {
  for i := 1; i<len(dungeon[0]) - 1; i++ {
    for j := 1; j<len(dungeon) - 1; j++ {
      if(dungeon[j-1][i-1] == WALL &&
        dungeon[j][i-1] == WALL &&
        dungeon[j+1][i-1] == WALL &&
        dungeon[j-1][i] == WALL &&
        dungeon[j][i] == WALL &&
        dungeon[j+1][i] == WALL &&
        dungeon[j-1][i+1] == WALL &&
        dungeon[j][i+1] == WALL &&
        dungeon[j+1][i+1] == WALL) {

      	// TODO: start new tunnel region
        continueMaze(dungeon, i, j)
      }
    }
  }
}

func continueMaze(dungeon [][]int, x int, y int) {
  validTiles := []Point{}

  if(x-2 >= 0 && dungeon[y][x-1] == WALL) {
    // check if is valid move by checking surroundings
    if(dungeon[y][x-2] == WALL &&
      dungeon[y+1][x-2] == WALL &&
      dungeon[y-1][x-2] == WALL &&
      dungeon[y+1][x-1] == WALL &&
      dungeon[y-1][x-1] == WALL) {

      validTiles = append(validTiles, Point{y: y, x: x-1})
    }
  }
  if(x+2 < dungeonWidth && dungeon[y][x+1] == WALL) {
    if(dungeon[y][x+2] == WALL &&
      dungeon[y-1][x+2] == WALL &&
      dungeon[y+1][x+2] == WALL &&
      dungeon[y+1][x+1] == WALL &&
      dungeon[y-1][x+1] == WALL) {

      validTiles = append(validTiles, Point{y: y, x: x+1})
    }
  }
  if(y-2 >= 0 && dungeon[y-1][x] == WALL) {
    if(dungeon[y-2][x] == WALL &&
      dungeon[y-2][x-1] == WALL &&
      dungeon[y-2][x+1] == WALL &&
      dungeon[y-1][x-1] == WALL &&
      dungeon[y-1][x+1] == WALL) {

      validTiles = append(validTiles, Point{y: y-1, x: x})
    }
  }
  if(y+2 < dungeonHeight && dungeon[y+1][x] == WALL) {
    if(dungeon[y+2][x] == WALL &&
      dungeon[y+2][x-1] == WALL &&
      dungeon[y+2][x+1] == WALL &&
      dungeon[y+1][x-1] == WALL &&
      dungeon[y+1][x+1] == WALL) {

      validTiles = append(validTiles, Point{y: y+1, x: x})
    }
  }

  if( len(validTiles) > 1 ) {
    i := rand.Intn( len(validTiles) )
    point := validTiles[i]
    dungeon[point.y][point.x] = TUNNEL
    continueMaze(dungeon, point.x, point.y)
    continueMaze(dungeon, x, y)
  } else if( len(validTiles) == 1 ) {
    point := validTiles[0]
    dungeon[point.y][point.x] = TUNNEL
    continueMaze(dungeon, point.x, point.y)
    continueMaze(dungeon, x, y)
  }
}

func identifyEdges(dungeon [][]int, rooms []Room) {
  for i := range rooms {
    x := rooms[i].x
    y := rooms[i].y

    for j := x; j < x + rooms[i].width; j++ {
      if(dungeon[y-2][j] == TUNNEL || dungeon[y-2][j] == FLOOR) {
        dungeon[y-1][j] = EDGE
        rooms[i].edges = append(rooms[i].edges, Point{ x: j, y: y-1})
      }
      if(dungeon[y+rooms[i].height+1][j] == TUNNEL || dungeon[y+rooms[i].height+1][j] == FLOOR) {
        dungeon[y+rooms[i].height][j] = EDGE
        rooms[i].edges = append(rooms[i].edges, Point{ x: j, y: y+rooms[i].height})
      }
    }

    for k := y; k < y + rooms[i].height; k++ {
      if(dungeon[k][x-2] == TUNNEL || dungeon[k][x-2] == FLOOR) {
        dungeon[k][x-1] = EDGE
        rooms[i].edges = append(rooms[i].edges, Point{ x: x-1, y: k})
      }
      if(dungeon[k][x+rooms[i].width+1] == TUNNEL || dungeon[k][x+rooms[i].width+1] == FLOOR) {
        dungeon[k][x+rooms[i].width] = EDGE
        rooms[i].edges = append(rooms[i].edges, Point{ x: x+rooms[i].width, y: k})
      }
    }
  }
}

func renderDungeon(dungeon [][]int) {
  for y := 0; y < dungeonHeight; y++ {
    for x := 0; x < dungeonWidth; x++ {
      switch dungeon[y][x] {
      case WALL:
        fmt.Print("0 ")
        break
      case FLOOR:
        fmt.Print("= ")
        break
      case EDGE:
        fmt.Print("* ")
        break
      case DOOR:
        fmt.Print("| ")
        break
      case TUNNEL:
        fmt.Print("- ")
        break
      default:
        fmt.Print("ER")
      }
    }

    fmt.Println()
  }
}

func main() {
  rand.Seed( time.Now().UTC().UnixNano())

  dungeon := createEmptyDungeon(dungeonWidth, dungeonHeight)
  rooms := createRooms(dungeon, minRoomSize, maxRoomSize, roomAttempts)
  createMaze(dungeon, rooms)
  identifyEdges(dungeon, rooms)
  // TODO: connect regions
  renderDungeon(dungeon)
}

