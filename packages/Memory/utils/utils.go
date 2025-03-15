package utils

import (
	"fmt"
	"main/packages/Memory/memory"
	"strings"
	"unsafe"
)

type Offsets struct {
	RenderViewFromRenderJob uint64
	DataModelHolder         uint64
	DataModel               uint64
	VisualDataModel         uint64

	Name            uint64
	Children        uint64
	Parent          uint64
	ClassDescriptor uint64
	LocalPlayer     uint64

	ValueBase   uint64
	ModuleFlags uint64
	IsCore      uint64

	PlaceID uint64

	BytecodeSize uint64
	Bytecode     map[string]uint64

	OffsetTaskScheduler uint64
	OffsetJobsContainer uint64
}

var OffsetsDataPlayer = Offsets{
	RenderViewFromRenderJob: 0x1E8, // Updated
	DataModelHolder:         0x118, // Updated
	DataModel:               0x1A8, // Updated
	VisualDataModel:         0x720, // Updated

	Name:            0x68,  // Updated
	Children:        0x70,  // Updated
	Parent:          0x50,  // Updated
	ClassDescriptor: 0x18,  // Updated
	ValueBase:       0xC8,  // Updated
	LocalPlayer:     0x120, // Updated

	ModuleFlags: 0x1B0 - 0x4, // Updated
	IsCore:      0x1B0,        // Updated

	PlaceID: 0x178, // Updated

	BytecodeSize: 0xA8, // Updated
	Bytecode: map[string]uint64{
		"LocalScript":  0x1C8, // Updated
		"ModuleScript": 0x170, // Updated
	},

	OffsetTaskScheduler: 0x61E5E38, // Updated
	OffsetJobsContainer: 0x1C8,     // Updated
}

var OffsetsDataUwp = Offsets{
	DataModelHolder: 0x118, // Updated
	DataModel:       0x1A8, // Updated
	VisualDataModel: 0x720, // Updated

	Name:            0x68,  // Updated
	Children:        0x70,  // Updated
	Parent:          0x50,  // Updated
	ClassDescriptor: 0x18,  // Updated
	ValueBase:       0xC8,  // Updated
	LocalPlayer:     0x120, // Updated

	ModuleFlags: 0x1B0 - 0x4, // Updated
	IsCore:      0x1B0,        // Updated

	BytecodeSize: 0xA8, // Updated
	Bytecode: map[string]uint64{
		"LocalScript":  0x1C8, // Updated
		"ModuleScript": 0x170, // Updated
	},
}

// GetDataModel retrieves the DataModel using the TaskScheduler and RenderJob
func GetDataModel(mem *memory.Luna, offsets Offsets) (uint64, error) {
	// Get the TaskScheduler pointer
	taskScheduler, err := mem.ReadPointer(mem.RobloxBase + uintptr(offsets.OffsetTaskScheduler))
	if err != nil {
		return 0, fmt.Errorf("failed to read TaskScheduler pointer: %v", err)
	}

	// Get the Jobs container
	jobsContainer, err := mem.ReadPointer(taskScheduler + uintptr(offsets.OffsetJobsContainer))
	if err != nil {
		return 0, fmt.Errorf("failed to read Jobs container: %v", err)
	}

	// Iterate through the jobs to find the RenderJob
	for i := 0; i < 0x500; i += 0x10 {
		job, err := mem.ReadPointer(jobsContainer + uintptr(i))
		if err != nil {
			continue
		}

		// Check if this is the RenderJob
		jobName, err := mem.ReadString(job+uintptr(offsets.Name), 10)
		if err != nil || !strings.EqualFold(jobName, "RenderJob") {
			continue
		}

		// Get the RenderView from the RenderJob
		renderView, err := mem.ReadPointer(job + uintptr(offsets.RenderViewFromRenderJob))
		if err != nil {
			return 0, fmt.Errorf("failed to read RenderView: %v", err)
		}

		// Get the FakeDataModel from the RenderJob
		fakeDataModel, err := mem.ReadPointer(job + uintptr(offsets.DataModelHolder))
		if err != nil {
			return 0, fmt.Errorf("failed to read FakeDataModel: %v", err)
		}

		// Get the real DataModel from the FakeDataModel
		dataModel, err := mem.ReadPointer(fakeDataModel + uintptr(offsets.DataModel))
		if err != nil {
			return 0, fmt.Errorf("failed to read DataModel: %v", err)
		}

		return uint64(dataModel), nil
	}

	return 0, fmt.Errorf("RenderJob not found")
}

// GetRenderVDM retrieves the RenderView and DataModel
func GetRenderVDM(pid uint32, mem *memory.Luna, offsets Offsets, UWP bool) uint64 {
	if mem == nil {
		return 0
	}

	// Get the DataModel
	dataModel, err := GetDataModel(mem, offsets)
	if err != nil {
		fmt.Println("Failed to get DataModel:", err)
		return 0
	}

	// Get the VisualDataModel from the DataModel
	var visualDataModel uint64
	err = mem.MemRead(uintptr(dataModel+offsets.VisualDataModel), unsafe.Pointer(&visualDataModel), unsafe.Sizeof(visualDataModel))
	if err != nil {
		fmt.Println("Failed to read VisualDataModel:", err)
		return 0
	}

	return visualDataModel
}
