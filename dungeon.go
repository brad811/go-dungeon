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

type Point struct {
  x int
  y int
}

type Tile struct {
  location Point
  region int
  material int
}

type Room struct {
  width int
  height int
  location Point
  edges []Point
}

type Dungeon struct {
  tiles [][]Tile
  rooms []Room
  width int
  height int
  numRegions int
}

func createEmptyDungeon(width int, height int) Dungeon {
  fmt.Println("Creating empty dungeon...")
  dungeon := Dungeon{ width: width, height: height }
  dungeon.tiles = make([][]Tile, height)
  for i := range dungeon.tiles {
    dungeon.tiles[i] = make([]Tile, width)
  }

  return dungeon
}

func createRooms(dungeon Dungeon, minSize, maxSize, attempts int) Dungeon {
  fmt.Println("Creating rooms...")
  var rooms []Room

  for i := 0; i < attempts; i++ {
    width := rand.Intn(maxSize - minSize) + minSize
    height := rand.Intn(maxSize - minSize) + minSize

    maxX := dungeon.width - width - 2
    maxY := dungeon.height - height - 2

    x := rand.Intn(maxX - 3) + 3
    y := rand.Intn(maxY - 3) + 3

    shouldAppend := true
    for r := range rooms {
      if(x + width < rooms[r].location.x || // to the left
      x > rooms[r].location.x + rooms[r].width || // to the right
      y + height < rooms[r].location.y || // fully above
      y > rooms[r].location.y + rooms[r].height) { // fully below
        // do nothing
      } else {
        shouldAppend = false
        break
      }
    }

    if(shouldAppend) {
      rooms = append(rooms, Room{ width: width, height: height, location: Point{ x: x, y: y } })
    }
  }

  for r := range rooms {
    dungeon.numRegions++
    for i := rooms[r].location.x; i < rooms[r].location.x + rooms[r].width; i++ {
      for j := rooms[r].location.y; j < rooms[r].location.y + rooms[r].height; j++ {
        dungeon.tiles[j][i].material = FLOOR
        dungeon.tiles[j][i].region = dungeon.numRegions
      }
    }
  }

  dungeon.rooms = rooms
  return dungeon
}

func createMaze(dungeon Dungeon) Dungeon {
  for i := 1; i<len(dungeon.tiles[0]) - 1; i++ {
    for j := 1; j<len(dungeon.tiles) - 1; j++ {
      if(dungeon.tiles[j-1][i-1].material == WALL &&
        dungeon.tiles[j][i-1].material == WALL &&
        dungeon.tiles[j+1][i-1].material == WALL &&
        dungeon.tiles[j-1][i].material == WALL &&
        dungeon.tiles[j][i].material == WALL &&
        dungeon.tiles[j+1][i].material == WALL &&
        dungeon.tiles[j-1][i+1].material == WALL &&
        dungeon.tiles[j][i+1].material == WALL &&
        dungeon.tiles[j+1][i+1].material == WALL) {

        dungeon.numRegions++
        continueMaze(dungeon, i, j)
      }
    }
  }

  return dungeon
}

func continueMaze(dungeon Dungeon, x int, y int) {
  validTiles := []Point{}

  if(x-2 >= 0 && dungeon.tiles[y][x-1].material == WALL) {
    // check if is valid move by checking surroundings
    if(dungeon.tiles[y][x-2].material == WALL &&
      dungeon.tiles[y+1][x-2].material == WALL &&
      dungeon.tiles[y-1][x-2].material == WALL &&
      dungeon.tiles[y+1][x-1].material == WALL &&
      dungeon.tiles[y-1][x-1].material == WALL) {

      validTiles = append(validTiles, Point{y: y, x: x-1})
    }
  }
  if(x+2 < dungeon.width && dungeon.tiles[y][x+1].material == WALL) {
    if(dungeon.tiles[y][x+2].material == WALL &&
      dungeon.tiles[y-1][x+2].material == WALL &&
      dungeon.tiles[y+1][x+2].material == WALL &&
      dungeon.tiles[y+1][x+1].material == WALL &&
      dungeon.tiles[y-1][x+1].material == WALL) {

      validTiles = append(validTiles, Point{y: y, x: x+1})
    }
  }
  if(y-2 >= 0 && dungeon.tiles[y-1][x].material == WALL) {
    if(dungeon.tiles[y-2][x].material == WALL &&
      dungeon.tiles[y-2][x-1].material == WALL &&
      dungeon.tiles[y-2][x+1].material == WALL &&
      dungeon.tiles[y-1][x-1].material == WALL &&
      dungeon.tiles[y-1][x+1].material == WALL) {

      validTiles = append(validTiles, Point{y: y-1, x: x})
    }
  }
  if(y+2 < dungeon.height && dungeon.tiles[y+1][x].material == WALL) {
    if(dungeon.tiles[y+2][x].material == WALL &&
      dungeon.tiles[y+2][x-1].material == WALL &&
      dungeon.tiles[y+2][x+1].material == WALL &&
      dungeon.tiles[y+1][x-1].material == WALL &&
      dungeon.tiles[y+1][x+1].material == WALL) {

      validTiles = append(validTiles, Point{y: y+1, x: x})
    }
  }

  if( len(validTiles) > 1 ) {
    i := rand.Intn( len(validTiles) )
    point := validTiles[i]
    dungeon.tiles[point.y][point.x].material = TUNNEL
    dungeon.tiles[point.y][point.x].region = dungeon.numRegions

    continueMaze(dungeon, point.x, point.y)
    continueMaze(dungeon, x, y)
  } else if( len(validTiles) == 1 ) {
    point := validTiles[0]
    dungeon.tiles[point.y][point.x].material = TUNNEL
    dungeon.tiles[point.y][point.x].region = dungeon.numRegions

    continueMaze(dungeon, point.x, point.y)
    continueMaze(dungeon, x, y)
  }
}

func identifyEdges(dungeon Dungeon) Dungeon {
  for i := range dungeon.rooms {
    x := dungeon.rooms[i].location.x
    y := dungeon.rooms[i].location.y

    for j := x; j < x + dungeon.rooms[i].width; j++ {
      if(dungeon.tiles[y-2][j].material == TUNNEL ||
        dungeon.tiles[y-2][j].material == FLOOR) {
        
        dungeon.tiles[y-1][j].material = EDGE
        dungeon.rooms[i].edges = append(dungeon.rooms[i].edges, Point{ x: j, y: y-1})
      }
      if(dungeon.tiles[y+dungeon.rooms[i].height+1][j].material == TUNNEL ||
        dungeon.tiles[y+dungeon.rooms[i].height+1][j].material == FLOOR) {
        
        dungeon.tiles[y+dungeon.rooms[i].height][j].material = EDGE
        dungeon.rooms[i].edges = append(dungeon.rooms[i].edges, Point{ x: j, y: y+dungeon.rooms[i].height})
      }
    }

    for k := y; k < y + dungeon.rooms[i].height; k++ {
      if(dungeon.tiles[k][x-2].material == TUNNEL ||
        dungeon.tiles[k][x-2].material == FLOOR) {
        
        dungeon.tiles[k][x-1].material = EDGE
        dungeon.rooms[i].edges = append(dungeon.rooms[i].edges, Point{ x: x-1, y: k})
      }
      if(dungeon.tiles[k][x+dungeon.rooms[i].width+1].material == TUNNEL ||
        dungeon.tiles[k][x+dungeon.rooms[i].width+1].material == FLOOR) {
        
        dungeon.tiles[k][x+dungeon.rooms[i].width].material = EDGE
        dungeon.rooms[i].edges = append(dungeon.rooms[i].edges, Point{ x: x+dungeon.rooms[i].width, y: k})
      }
    }
  }

  return dungeon
}

func connectRegions(dungeon Dungeon) Dungeon {
  for i := range dungeon.rooms {
    room := dungeon.rooms[i]
    edge := room.edges[ rand.Intn( len(dungeon.rooms[i].edges) ) ]
    roomRegion := dungeon.tiles[ dungeon.rooms[i].location.y ][ dungeon.rooms[i].location.x ].region

    // check if edge is unconnected
    surroundingTiles := [8]Tile{}
    surroundingTiles[0] = dungeon.tiles[edge.y-1][edge.x-1]
    surroundingTiles[1] = dungeon.tiles[edge.y-1][edge.x]
    surroundingTiles[2] = dungeon.tiles[edge.y-1][edge.x+1]
    surroundingTiles[3] = dungeon.tiles[edge.y][edge.x-1]
    surroundingTiles[4] = dungeon.tiles[edge.y][edge.x+1]
    surroundingTiles[5] = dungeon.tiles[edge.y+1][edge.x-1]
    surroundingTiles[6] = dungeon.tiles[edge.y+1][edge.x]
    surroundingTiles[7] = dungeon.tiles[edge.y+1][edge.x+1]

    for j := range surroundingTiles {
      if((surroundingTiles[j].material == FLOOR || surroundingTiles[j].material == TUNNEL) &&
        surroundingTiles[j].region != roomRegion) {

        dungeon.tiles[edge.y][edge.x].material = DOOR
        for x := room.location.x-1; x < room.location.x + room.width - 1; x++ {
          for y := room.location.y-1; y < room.location.y + room.height - 1; y++ {
            dungeon.tiles[y][x].region = surroundingTiles[j].region
          }
        }

        break
      }
    }
  }

  return dungeon
}

func renderDungeon(dungeon Dungeon) {
  fmt.Println("Dungeon: (",dungeon.width,",",dungeon.height,") Regions: ", dungeon.numRegions)

  for y := 0; y < dungeon.height; y++ {
    for x := 0; x < dungeon.width; x++ {
      switch dungeon.tiles[y][x].material {
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
  dungeon = createRooms(dungeon, minRoomSize, maxRoomSize, roomAttempts)
  dungeon = createMaze(dungeon)
  dungeon = identifyEdges(dungeon)
  dungeon = connectRegions(dungeon)
  renderDungeon(dungeon)
}

