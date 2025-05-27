package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/GeoNet/delta/meta"
)

type DB struct {
	db *sql.DB
}

func New(db *sql.DB) DB {
	return DB{
		db: db,
	}
}

func (d DB) exec(ctx context.Context, tx *sql.Tx, cmds ...string) error {
	for _, cmd := range cmds {
		if _, err := tx.ExecContext(ctx, cmd); err != nil {
			return fmt.Errorf("cmd %q: %w", cmd, err)
		}
	}
	return nil
}

func (d DB) prepare(ctx context.Context, tx *sql.Tx, cmd string, values ...[]any) error {

	stmt, err := tx.PrepareContext(ctx, cmd)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, v := range values {
		if _, err := stmt.ExecContext(ctx, v...); err != nil {
			return fmt.Errorf("%v : %w", v, err)
		}
	}

	return nil
}

func (d DB) Init(ctx context.Context, list []meta.TableList) error {

	tables := make(map[string]meta.TableList)
	for _, v := range list {
		tables[v.Table.Name()] = v
	}

	// Get a Tx for making transaction requests.
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	// Defer a rollback in case anything fails, not actually
	// worried about any rollback error.
	defer func() { _ = tx.Rollback() }()

	// overall lookup tables, this should likely be pre-populated
	// to allow for constraints to be applied.
	if err := d.exec(ctx, tx, datum.Create); err != nil {
		return err
	}

	if err := d.exec(ctx, tx, reference.Create); err != nil {
		return fmt.Errorf("reference create: %v", err)
	}
	if err := d.exec(ctx, tx, referenceNetwork.Create); err != nil {
		return fmt.Errorf("reference network create: %v", err)
	}

	if err := d.exec(ctx, tx, method.Create); err != nil {
		return fmt.Errorf("method create: %v", err)
	}
	if err := d.exec(ctx, tx, placeRole.Create); err != nil {
		return fmt.Errorf("place role create: %v", err)
	}

	for _, l := range list {
		switch l.Table.Name() {
		case "Response":
			if err := d.exec(ctx, tx, response.Create); err != nil {
				return fmt.Errorf("resonse create: %v", err)
			}
			if err := d.prepare(ctx, tx, response.Insert(), response.Columns(l)...); err != nil {
				return fmt.Errorf("response insert: %v", err)
			}
		case "Placename":
			if err := d.exec(ctx, tx, placename.Create); err != nil {
				return fmt.Errorf("placename create: %v", err)
			}
			if err := d.prepare(ctx, tx, placename.Insert(), placename.Columns(l)...); err != nil {
				return fmt.Errorf("placename insert: %v", err)
			}
		case "Citation":
			if err := d.exec(ctx, tx, citation.Create); err != nil {
				return fmt.Errorf("citation create: %v", err)
			}
			if err := d.prepare(ctx, tx, citation.Insert(), citation.Columns(l)...); err != nil {
				return fmt.Errorf("citation insert: %v", err)
			}
		case "Asset":
			if err := d.exec(ctx, tx, makeCreate); err != nil {
				return fmt.Errorf("make create: %v", err)
			}
			if err := d.prepare(ctx, tx, mmake.Insert(), mmake.Columns(l)...); err != nil {
				return fmt.Errorf("make insert: %v", err)
			}
			if err := d.exec(ctx, tx, modelCreate); err != nil {
				return fmt.Errorf("model create: %v", err)
			}
			if err := d.prepare(ctx, tx, model.Insert(), model.Columns(l)...); err != nil {
				return fmt.Errorf("model insert: %v", err)
			}
			if err := d.exec(ctx, tx, asset.Create); err != nil {
				return fmt.Errorf("asset create: %v", err)
			}
			if err := d.prepare(ctx, tx, asset.Insert(), asset.Columns(l)...); err != nil {
				return fmt.Errorf("asset insert: %v", err)
			}
		case "Firmware":
			if err := d.exec(ctx, tx, firmware.Create); err != nil {
				return fmt.Errorf("firmware create: %v", err)
			}
			if err := d.prepare(ctx, tx, firmware.Insert(), firmware.Columns(l)...); err != nil {
				return fmt.Errorf("firmware insert: %v", err)
			}
		case "Calibration":
			if err := d.exec(ctx, tx, calibration.Create); err != nil {
				return fmt.Errorf("calibration create: %v", err)
			}
			if err := d.prepare(ctx, tx, calibration.Insert(), calibration.Columns(l)...); err != nil {
				return fmt.Errorf("calibration insert: %v", err)
			}
		case "Channel":
			if err := d.prepare(ctx, tx, mmake.Insert(), mmake.Columns(l)...); err != nil {
				return fmt.Errorf("channel make insert: %v", err)
			}
			if err := d.prepare(ctx, tx, model.Insert(), model.Columns(l)...); err != nil {
				return fmt.Errorf("channel model insert: %v", err)
			}
			if err := d.exec(ctx, tx, channel.Create); err != nil {
				return fmt.Errorf("channel create: %v", err)
			}
			if err := d.prepare(ctx, tx, channel.Insert(), channel.Columns(l)...); err != nil {
				return fmt.Errorf("channel insert: %v", err)
			}
		case "Component":
			if err := d.prepare(ctx, tx, mmake.Insert(), mmake.Columns(l)...); err != nil {
				return fmt.Errorf("component make insert: %v", err)
			}
			if err := d.prepare(ctx, tx, model.Insert(), model.Columns(l)...); err != nil {
				return fmt.Errorf("component model insert: %v", err)
			}
			if err := d.exec(ctx, tx, component.Create); err != nil {
				return fmt.Errorf("component create: %v", err)
			}
			if err := d.prepare(ctx, tx, component.Insert(), component.Columns(l)...); err != nil {
				return fmt.Errorf("component insert: %v", err)
			}
		case "Network":
			if err := d.exec(ctx, tx, network.Create); err != nil {
				return fmt.Errorf("network create: %v", err)
			}
			if err := d.prepare(ctx, tx, network.Insert(), network.Columns(l)...); err != nil {
				return fmt.Errorf("network insert: %v", err)
			}
		case "Station":
			if err := d.prepare(ctx, tx, datum.Insert(), datum.Columns(l)...); err != nil {
				return fmt.Errorf("datum insert: %v", err)
			}
			if err := d.exec(ctx, tx, station.Create); err != nil {
				return fmt.Errorf("station create: %v", err)
			}
			if err := d.prepare(ctx, tx, station.Insert(), station.Columns(l)...); err != nil {
				return fmt.Errorf("station insert: %v", err)
			}
			if err := d.exec(ctx, tx, stationNetwork.Create); err != nil {
				return fmt.Errorf("station network create: %v", err)
			}
			if err := d.prepare(ctx, tx, stationNetwork.Insert(), stationNetwork.Columns(l)...); err != nil {
				return fmt.Errorf("station network insert: %v", err)
			}
		case "Sample":
			if err := d.prepare(ctx, tx, datum.Insert(), datum.Columns(l)...); err != nil {
				return fmt.Errorf("datum insert: %v", err)
			}
			if err := d.exec(ctx, tx, sample.Create); err != nil {
				return fmt.Errorf("sample create: %v", err)
			}
			if err := d.prepare(ctx, tx, sample.Insert(), sample.Columns(l)...); err != nil {
				return fmt.Errorf("sample insert: %v", err)
			}
			if err := d.exec(ctx, tx, sampleNetwork.Create); err != nil {
				return fmt.Errorf("sample network create: %v", err)
			}
			if err := d.prepare(ctx, tx, sampleNetwork.Insert(), sampleNetwork.Columns(l)...); err != nil {
				return fmt.Errorf("sample network insert: %v", err)
			}
		case "Site":
			if err := d.prepare(ctx, tx, datum.Insert(), datum.Columns(l)...); err != nil {
				return fmt.Errorf("datum insert: %v", err)
			}
			if err := d.exec(ctx, tx, site.Create); err != nil {
				return fmt.Errorf("site create: %v", err)
			}
			if err := d.prepare(ctx, tx, site.Insert(), site.Columns(l)...); err != nil {
				return fmt.Errorf("site insert: %v", err)
			}
		case "Point":
			if err := d.prepare(ctx, tx, datum.Insert(), datum.Columns(l)...); err != nil {
				return fmt.Errorf("datum insert: %v", err)
			}
			if err := d.exec(ctx, tx, point.Create); err != nil {
				return fmt.Errorf("point create: %v", err)
			}
			if err := d.prepare(ctx, tx, point.Insert(), point.Columns(l)...); err != nil {
				return fmt.Errorf("point insert: %v", err)
			}
		case "Feature":
			if err := d.exec(ctx, tx, feature.Create); err != nil {
				return fmt.Errorf("feature create: %v", err)
			}
			if err := d.prepare(ctx, tx, feature.Insert(), feature.Columns(l)...); err != nil {
				return fmt.Errorf("feature insert: %v", err)
			}
		case "Class":
			if err := d.exec(ctx, tx, class.Create); err != nil {
				return fmt.Errorf("class create: %v", err)
			}
			if err := d.prepare(ctx, tx, class.Insert(), class.Columns(l)...); err != nil {
				return fmt.Errorf("class insert: %v", err)
			}
			if err := d.exec(ctx, tx, classCitation.Create); err != nil {
				return fmt.Errorf("class citation create: %v", err)
			}
			if err := d.prepare(ctx, tx, classCitation.Insert(), class.Links(l, "Station")...); err != nil {
				return fmt.Errorf("class citation insert: %v", err)
			}
		case "Mark":
			if err := d.exec(ctx, tx, mark.Create); err != nil {
				return fmt.Errorf("mark create: %v", err)
			}
			if err := d.prepare(ctx, tx, mark.Insert(), mark.Columns(l)...); err != nil {
				return fmt.Errorf("mark insert: %v", err)
			}
			if err := d.exec(ctx, tx, markNetwork.Create); err != nil {
				return fmt.Errorf("mark network create: %v", err)
			}
			if err := d.prepare(ctx, tx, markNetwork.Insert(), markNetwork.Columns(l)...); err != nil {
				return fmt.Errorf("mark network insert: %v", err)
			}
		case "Monument":
			if err := d.exec(ctx, tx, markType.Create); err != nil {
				return fmt.Errorf("mark type create: %v", err)
			}
			if err := d.prepare(ctx, tx, markType.Insert(), markType.Columns(l)...); err != nil {
				return fmt.Errorf("mark type insert: %v", err)
			}
			if err := d.exec(ctx, tx, monumentType.Create); err != nil {
				return fmt.Errorf("monument type create: %v", err)
			}
			if err := d.prepare(ctx, tx, monumentType.Insert(), monumentType.Columns(l)...); err != nil {
				return fmt.Errorf("monument type insert: %v", err)
			}
			if err := d.exec(ctx, tx, foundationType.Create); err != nil {
				return fmt.Errorf("foundation type create: %v", err)
			}
			if err := d.prepare(ctx, tx, foundationType.Insert(), foundationType.Columns(l)...); err != nil {
				return fmt.Errorf("foundation type insert: %v", err)
			}
			if err := d.exec(ctx, tx, bedrock.Create); err != nil {
				return fmt.Errorf("bedrock create: %v", err)
			}
			if err := d.prepare(ctx, tx, bedrock.Insert(), bedrock.Columns(l)...); err != nil {
				return fmt.Errorf("bedrock insert: %v", err)
			}
			if err := d.exec(ctx, tx, geology.Create); err != nil {
				return fmt.Errorf("geology create: %v", err)
			}
			if err := d.prepare(ctx, tx, geology.Insert(), geology.Columns(l)...); err != nil {
				return fmt.Errorf("geology insert: %v", err)
			}
			if err := d.exec(ctx, tx, monument.Create); err != nil {
				return fmt.Errorf("monument create: %v", err)
			}
			if err := d.prepare(ctx, tx, monument.Insert(), monument.Columns(l)...); err != nil {
				return fmt.Errorf("monument insert: %v", err)
			}
		case "Visibility":
			if err := d.exec(ctx, tx, visibility.Create); err != nil {
				return fmt.Errorf("visibility create: %v", err)
			}
			if err := d.prepare(ctx, tx, visibility.Insert(), visibility.Columns(l)...); err != nil {
				return fmt.Errorf("visibility insert: %v", err)
			}
		case "Antenna":
			if err := d.exec(ctx, tx, antenna.Create); err != nil {
				return fmt.Errorf("antenna create: %v", err)
			}
			if err := d.prepare(ctx, tx, antenna.Insert(), antenna.Columns(l)...); err != nil {
				return fmt.Errorf("antenna insert: %v", err)
			}
		case "MetSensor":
			if err := d.exec(ctx, tx, metsensor.Create); err != nil {
				return fmt.Errorf("metsensor create: %v", err)
			}
			if err := d.prepare(ctx, tx, metsensor.Insert(), metsensor.Columns(l)...); err != nil {
				return fmt.Errorf("metsensor insert: %v", err)
			}
		case "Radome":
			if err := d.exec(ctx, tx, radome.Create); err != nil {
				return fmt.Errorf("radome create: %v", err)
			}
			if err := d.prepare(ctx, tx, radome.Insert(), radome.Columns(l)...); err != nil {
				return fmt.Errorf("radome insert: %v", err)
			}
		case "Receiver":
			if err := d.exec(ctx, tx, receiver.Create); err != nil {
				return fmt.Errorf("receiver create: %v", err)
			}
			if err := d.prepare(ctx, tx, receiver.Insert(), receiver.Columns(l)...); err != nil {
				return fmt.Errorf("receiver insert: %v", err)
			}
		case "Session":
			if err := d.exec(ctx, tx, session.Create); err != nil {
				return fmt.Errorf("session create: %v", err)
			}
			if err := d.prepare(ctx, tx, session.Insert(), session.Columns(l)...); err != nil {
				return fmt.Errorf("session insert: %v", err)
			}
		case "Gauge":
			if err := d.exec(ctx, tx, gauge.Create); err != nil {
				return fmt.Errorf("gauge create: %v", err)
			}
			if err := d.prepare(ctx, tx, gauge.Insert(), gauge.Columns(l)...); err != nil {
				return fmt.Errorf("gauge insert: %v", err)
			}
		case "Constituent":
			if err := d.exec(ctx, tx, constituent.Create); err != nil {
				return fmt.Errorf("constituent create: %v", err)
			}
			if err := d.prepare(ctx, tx, constituent.Insert(), constituent.Columns(l)...); err != nil {
				return fmt.Errorf("constituent insert: %v", err)
			}
		case "Dart":
			if err := d.exec(ctx, tx, dartCreate); err != nil {
				return fmt.Errorf("dart create: %v", err)
			}
			if err := d.prepare(ctx, tx, dart.Insert(), dart.Columns(l)...); err != nil {
				return fmt.Errorf("dart insert: %v", err)
			}
		case "Mount":
			if err := d.prepare(ctx, tx, datum.Insert(), datum.Columns(l)...); err != nil {
				return fmt.Errorf("datum insert: %v", err)
			}
			if err := d.exec(ctx, tx, mount.Create); err != nil {
				return fmt.Errorf("mount create: %v", err)
			}
			if err := d.prepare(ctx, tx, mount.Insert(), mount.Columns(l)...); err != nil {
				return fmt.Errorf("mount insert: %v", err)
			}
			if err := d.exec(ctx, tx, mountNetwork.Create); err != nil {
				return fmt.Errorf("mount network create: %v", err)
			}
			if err := d.prepare(ctx, tx, mountNetwork.Insert(), mountNetwork.Columns(l)...); err != nil {
				return fmt.Errorf("mount network insert: %v", err)
			}
		case "View":
			if err := d.exec(ctx, tx, view.Create); err != nil {
				return fmt.Errorf("view create: %v", err)
			}
			if err := d.prepare(ctx, tx, view.Insert(), view.Columns(l)...); err != nil {
				return fmt.Errorf("view insert: %v", err)
			}
		case "Camera":
			if err := d.exec(ctx, tx, camera.Create); err != nil {
				return fmt.Errorf("camera create: %v", err)
			}
			if err := d.prepare(ctx, tx, camera.Insert(), camera.Columns(l)...); err != nil {
				return fmt.Errorf("camera insert: %v", err)
			}
		case "Doas":
			if err := d.exec(ctx, tx, doas.Create); err != nil {
				return fmt.Errorf("doas create: %v", err)
			}
			if err := d.prepare(ctx, tx, doas.Insert(), doas.Columns(l)...); err != nil {
				return fmt.Errorf("doas insert: %v", err)
			}
		case "Datalogger":
			if err := d.prepare(ctx, tx, placeRole.Insert(), placeRole.Columns(l)...); err != nil {
				return fmt.Errorf("place role insert: %v", err)
			}
			if err := d.exec(ctx, tx, datalogger.Create); err != nil {
				return fmt.Errorf("datalogger create: %v", err)
			}
			if err := d.prepare(ctx, tx, datalogger.Insert(), datalogger.Columns(l)...); err != nil {
				return fmt.Errorf("datalogger insert: %v", err)
			}
		case "Sensor":
			if err := d.prepare(ctx, tx, method.Insert(), method.Columns(l)...); err != nil {
				return fmt.Errorf("method insert: %v", err)
			}
			if err := d.exec(ctx, tx, sensor.Create); err != nil {
				return fmt.Errorf("sensor create: %v", err)
			}
			if err := d.prepare(ctx, tx, sensor.Insert(), sensor.Columns(l)...); err != nil {
				return fmt.Errorf("sensor insert: %v", err)
			}
		case "Recorder":
			if err := d.prepare(ctx, tx, method.Insert(), method.Columns(l)...); err != nil {
				return fmt.Errorf("method insert: %v", err)
			}
			if err := d.exec(ctx, tx, recorder.Create); err != nil {
				return fmt.Errorf("method create: %v", err)
			}
			if err := d.prepare(ctx, tx, recorderModel.Insert(), recorderModel.Columns(l)...); err != nil {
				return fmt.Errorf("recorder model insert: %v", err)
			}
			if err := d.prepare(ctx, tx, recorder.Insert(), recorder.Columns(l)...); err != nil {
				return fmt.Errorf("recorder insert: %v", err)
			}
		case "Timing":
			if err := d.exec(ctx, tx, timing.Create); err != nil {
				return fmt.Errorf("timing create: %v", err)
			}
			if err := d.prepare(ctx, tx, timing.Insert(), timing.Columns(l)...); err != nil {
				return fmt.Errorf("timing insert: %v", err)
			}
		case "Telemetry":
			if err := d.exec(ctx, tx, telemetry.Create); err != nil {
				return fmt.Errorf("telemetry create: %v", err)
			}
			if err := d.prepare(ctx, tx, telemetry.Insert(), telemetry.Columns(l)...); err != nil {
				return fmt.Errorf("telemetry insert: %v", err)
			}
		case "Polarity":
			if err := d.exec(ctx, tx, polarity.Create); err != nil {
				return fmt.Errorf("polarity create: %v", err)
			}
			if err := d.prepare(ctx, tx, polarity.Insert(), polarity.Columns(l)...); err != nil {
				return fmt.Errorf("polarity insert: %v", err)
			}
		case "Gain":
			if err := d.exec(ctx, tx, gain.Create); err != nil {
				return fmt.Errorf("gain create: %v", err)
			}
			if err := d.prepare(ctx, tx, gain.Insert(), gain.Columns(l)...); err != nil {
				return fmt.Errorf("gain insert: %v", err)
			}
		case "Preamp":
			if err := d.exec(ctx, tx, preamp.Create); err != nil {
				return fmt.Errorf("preamp create: %v", err)
			}
			if err := d.prepare(ctx, tx, preamp.Insert(), preamp.Columns(l)...); err != nil {
				return fmt.Errorf("preamp insert: %v", err)
			}
		case "Stream":
			if err := d.exec(ctx, tx, stream.Create); err != nil {
				return fmt.Errorf("stream create: %v", err)
			}
			if err := d.prepare(ctx, tx, stream.Insert(), stream.Columns(l)...); err != nil {
				return fmt.Errorf("stream insert: %v", err)
			}
		case "Connection":
			if err := d.prepare(ctx, tx, placeRole.Insert(), placeRole.Columns(l)...); err != nil {
				return fmt.Errorf("connection place role insert: %v", err)
			}
			if err := d.exec(ctx, tx, connection.Create); err != nil {
				return fmt.Errorf("connection create: %v", err)
			}
			if err := d.prepare(ctx, tx, connection.Insert(), connection.Columns(l)...); err != nil {
				return fmt.Errorf("connection insert: %v", err)
			}
		default:
			log.Printf("ignoring %s", l.Table.Name())
		}
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
