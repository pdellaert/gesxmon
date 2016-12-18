package vsphere

import (
	"context"
	"net/url"
	"reflect"

	"github.com/Sirupsen/logrus"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/event"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/vim25/types"
)

// ESXiEventListener represents a listener for the events we want to capture from ESXi
type ESXiEventListener struct {
	client   *govmomi.Client
	url      *url.URL
	insecure bool
	logger   *logrus.Logger
}

// NewESXiEventListener returns a new ESXiEventListener
func NewESXiEventListener(url *url.URL, insecure bool, logger *logrus.Logger) *ESXiEventListener {
	return &ESXiEventListener{
		url:      url,
		insecure: insecure,
		logger:   logger,
	}
}

// Run will run the listener and make it listen to the appropriate events
func (e *ESXiEventListener) Run(ctx context.Context) (err error) {
	e.logger.Debug("Connecting to the ESXi host")
	e.client, err = govmomi.NewClient(ctx, e.url, e.insecure)
	if err != nil {
		return err
	}
	e.logger.Info("Successful connected to the ESXi host")

	e.logger.Debug("Getting default DC for the ESXi host")
	finder := find.NewFinder(e.client.Client, true)
	dc, err := finder.DefaultDatacenter(ctx)
	if err != nil {
		return err
	}
	refs := []types.ManagedObjectReference{dc.Reference()}

	e.logger.Debug("Setting up the event manager")
	eventManager := event.NewManager(e.client.Client)
	err = eventManager.Events(ctx, refs, 10, true, false, e.handleEvent)
	if err != nil {
		return err
	}
	e.logger.Info("Successfully created the event manager")

	return nil
}

func (e *ESXiEventListener) handleEvent(ref types.ManagedObjectReference, events []types.BaseEvent) (err error) {

	for _, event := range events {
		switch event.(type) {
		case *types.VmBeingCreatedEvent:
			vmName := event.GetEvent().Vm.Name
			vmRef := event.GetEvent().Vm.Vm.Reference()
			e.logger.WithFields(logrus.Fields{"vm-name": vmName, "vm-ref": vmRef}).Info("VM being created event received")
		case *types.VmCreatedEvent:
			vmName := event.GetEvent().Vm.Name
			vmRef := event.GetEvent().Vm.Vm.Reference()
			e.logger.WithFields(logrus.Fields{"vm-name": vmName, "vm-ref": vmRef}).Info("VM created event received")
		case *types.VmRemovedEvent:
			vmName := event.GetEvent().Vm.Name
			vmRef := event.GetEvent().Vm.Vm.Reference()
			e.logger.WithFields(logrus.Fields{"vm-name": vmName, "vm-ref": vmRef}).Info("VM removed event received")
		case *types.VmStartingEvent:
			vmName := event.GetEvent().Vm.Name
			vmRef := event.GetEvent().Vm.Vm.Reference()
			e.logger.WithFields(logrus.Fields{"vm-name": vmName, "vm-ref": vmRef}).Info("VM starting event received")
		case *types.VmPoweredOnEvent:
			vmName := event.GetEvent().Vm.Name
			vmRef := event.GetEvent().Vm.Vm.Reference()
			e.logger.WithFields(logrus.Fields{"vm-name": vmName, "vm-ref": vmRef}).Info("VM powered on event received")
		case *types.VmSuspendingEvent:
			vmName := event.GetEvent().Vm.Name
			vmRef := event.GetEvent().Vm.Vm.Reference()
			e.logger.WithFields(logrus.Fields{"vm-name": vmName, "vm-ref": vmRef}).Info("VM suspending event received")
		case *types.VmSuspendedEvent:
			vmName := event.GetEvent().Vm.Name
			vmRef := event.GetEvent().Vm.Vm.Reference()
			e.logger.WithFields(logrus.Fields{"vm-name": vmName, "vm-ref": vmRef}).Info("VM suspended event received")
		case *types.VmResumingEvent:
			vmName := event.GetEvent().Vm.Name
			vmRef := event.GetEvent().Vm.Vm.Reference()
			e.logger.WithFields(logrus.Fields{"vm-name": vmName, "vm-ref": vmRef}).Info("VM resuming event received")
		case *types.VmStoppingEvent:
			vmName := event.GetEvent().Vm.Name
			vmRef := event.GetEvent().Vm.Vm.Reference()
			e.logger.WithFields(logrus.Fields{"vm-name": vmName, "vm-ref": vmRef}).Info("VM stopping event received")
		case *types.VmPoweredOffEvent:
			vmName := event.GetEvent().Vm.Name
			vmRef := event.GetEvent().Vm.Vm.Reference()
			e.logger.WithFields(logrus.Fields{"vm-name": vmName, "vm-ref": vmRef}).Info("VM powered off event received")
		case *types.VmResettingEvent:
			vmName := event.GetEvent().Vm.Name
			vmRef := event.GetEvent().Vm.Vm.Reference()
			e.logger.WithFields(logrus.Fields{"vm-name": vmName, "vm-ref": vmRef}).Info("VM resetting event received")
		case *types.VmRegisteredEvent:
			vmName := event.GetEvent().Vm.Name
			vmRef := event.GetEvent().Vm.Vm.Reference()
			e.logger.WithFields(logrus.Fields{"vm-name": vmName, "vm-ref": vmRef}).Info("VM registered event received")
		case *types.VmReconfiguredEvent:
			vmName := event.GetEvent().Vm.Name
			vmRef := event.GetEvent().Vm.Vm.Reference()
			e.logger.WithFields(logrus.Fields{"vm-name": vmName, "vm-ref": vmRef}).Info("VM reconfigure event received")
		default:
			e.logger.WithFields(logrus.Fields{"type": reflect.TypeOf(event).String()}).Debug("Event ignored")
		}
	}
	return nil
}
