[![Go Reference](https://pkg.go.dev/badge/github.com/padilin/gengeno.svg)](https://pkg.go.dev/github.com/padilin/gengeno) [![Go Report Card](https://goreportcard.com/badge/github.com/padilin/gengeno)](https://goreportcard.com/report/github.com/padilin/gengeno)
![Aseprite](https://img.shields.io/badge/Aseprite-FFFFFF?style=for-the-badge&logo=Aseprite&logoColor=#7D929E)
# Gengeno

A generator game written in Go using Ebitengine.

## Goal
Other than learning Go and Ebitengine, I want to create a game about generating power employing anything from coal, fuel, hydro, nucelar, and more.  The goal is to create semi-realistic simulation of power generation and distribution.

## Ideas Implemented (in no particular order)

- Pipes and Reservoirs
- Pressure based flow
- Sprites
- Sprites updated based on state
- Isometric view

## Ideas Not Implemented (in no particular order)

- Generators
- Power simulation
- Material properties
- Heat simulation
- UI
- Non-placeholder assets

## Program Flow

```mermaid
graph TD
    subgraph "Main Loop (Ebiten)"
        Main["main.go"] -->|RunGame| Game
        Game -->|Update| System
        Game -->|Draw| Level
    end

    subgraph "Data Structures"
        Game["Game Struct"]
        System["System Struct"]
        Level["Level Struct"]
        
        Game --> System
        Game --> Level
        
        System -->|Manages| NodeList["Nodes []Component"]
        System -->|Manages| PipeList["Pipes []*Pipe"]
        
        Level -->|Contains| Entities["Entities []*Entity"]
        Level -->|Contains| Tiles["Tiles [][]*Tile"]
        
        Entity -->|Has A| Component
        Entity -->|Has A| Sprite
        
        Pipe -->|Implements| Component
        Reservoir -->|Implements| Component
        Generator -->|Implements| Component
    end

    subgraph "Simulation Logic (System.Tick)"
        MethodTick["Tick()"]
        CalcFlow[Calculate Bernoulli Flow]
        MoveVol[Move Volume]
        UpdateNodes[Update Component States]
        
        System --> MethodTick
        MethodTick -->|Iterate Pipes| CalcFlow
        CalcFlow -->|Delta Head & Friction| MoveVol
        MoveVol -->|Update Pending| PipeList
        MethodTick -->|Apply Pending Changes| UpdateNodes
        UpdateNodes --> NodeList
    end

    subgraph "Rendering (Game.Draw)"
        MethodDraw["Draw()"]
        RenderLevel[renderLevel]
        DrawEntity["Entity.Draw"]
        
        Game --> MethodDraw
        MethodDraw --> RenderLevel
        RenderLevel -->|Iterate| Tiles
        Tiles -->|Contains| Entity
        Entity --> DrawEntity
        DrawEntity -->|Uses| Sprite
    end
```
