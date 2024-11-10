package main

import (
	"fmt"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/hierarchy"
)

var Pos = donburi.NewComponentType[PosData]()
var Parent = donburi.NewTag().SetName("parent")
var Child = donburi.NewTag().SetName("child")

type PosData struct { X, Y int }

func main() {
    w := donburi.NewWorld()
	ecsRef := ecs.NewECS(w)

	// ecsRef.AddSystem(hierarchy.HierarchySystem.RemoveChildren)
	ecsRef.AddSystem(func(ecs *ecs.ECS) {
        Child.Each(ecs.World, func(e *donburi.Entry){
            fmt.Println("child", ecs.World.Valid(e.Entity()))
        })
    })

    pe := w.Entry(w.Create(Parent))
	ce1 := w.Entry(w.Create(Child, Pos))
	ce2 := w.Entry(w.Create(Child, Pos))
    
    Pos.Set(ce1, &PosData{5, 10})
    Pos.Set(ce2, &PosData{100, 69})

    // hierarchy.SetParent(ce1, pe)
    // hierarchy.SetParent(ce2, pe)
    hierarchy.AppendChild(pe, ce1)
    hierarchy.AppendChild(pe, ce2)

    hierarchy.RemoveChildrenRecursive(pe)

    // entry := w.Entry(w.Create(Pos))
    // Pos.Set(entry, &PosData{101, 1010})

    ecsRef.Update()

    fmt.Println("has", hierarchy.HasChildren(pe))
    children, ok := hierarchy.GetChildren(pe)
    // fmt.Println("old child", ce1)
    fmt.Println("children", children, ok)

    Pos.Each(w, func(e *donburi.Entry) {
        fmt.Println(">", Pos.Get(e))
    })
}