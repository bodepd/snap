/*
http://www.apache.org/licenses/LICENSE-2.0.txt


Copyright 2015 Intel Corporation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package plugin

import (
	"strings"
	"testing"
	"time"

	"github.com/intelsdi-x/snap/core/cdata"
	"github.com/intelsdi-x/snap/core/ctypes"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMetric(t *testing.T) {
	Convey("error on invalid snap content type", t, func() {
		m := []PluginMetricType{
			*NewPluginMetricType([]string{"foo", "bar"}, time.Now(), "", nil, nil, 1),
			*NewPluginMetricType([]string{"foo", "baz"}, time.Now(), "", nil, nil, 2),
		}
		a, c, e := MarshalPluginMetricTypes("foo", m)
		m[0].Version_ = 1
		m[0].AddData(3)
		configNewNode := cdata.NewNode()
		configNewNode.AddItem("user", ctypes.ConfigValueStr{Value: "foo"})
		m[0].Config_ = configNewNode
		So(e, ShouldNotBeNil)
		So(e.Error(), ShouldEqual, "invalid snap content type: foo")
		So(a, ShouldBeNil)
		So(c, ShouldEqual, "")
		So(m[0].Version(), ShouldResemble, 1)
		So(m[0].Data(), ShouldResemble, 3)
		So(m[0].Config(), ShouldNotBeNil)
	})

	Convey("error on empty metric slice", t, func() {
		m := []PluginMetricType{}
		a, c, e := MarshalPluginMetricTypes("foo", m)
		So(e, ShouldNotBeNil)
		So(e.Error(), ShouldEqual, "attempt to marshall empty slice of metrics: foo")
		So(a, ShouldBeNil)
		So(c, ShouldEqual, "")
	})

	Convey("marshall using snap.* default to snap.gob", t, func() {
		m := []PluginMetricType{
			*NewPluginMetricType([]string{"foo", "bar"}, time.Now(), "", nil, nil, 1),
			*NewPluginMetricType([]string{"foo", "baz"}, time.Now(), "", nil, nil, "2"),
		}
		a, c, e := MarshalPluginMetricTypes("snap.*", m)
		So(e, ShouldBeNil)
		So(a, ShouldNotBeNil)
		So(len(a), ShouldBeGreaterThan, 0)
		So(c, ShouldEqual, "snap.gob")

		Convey("unmarshal snap.gob", func() {
			m, e = UnmarshallPluginMetricTypes("snap.gob", a)
			So(e, ShouldBeNil)
			So(strings.Join(m[0].Namespace(), "/"), ShouldResemble, "foo/bar")
			So(m[0].Data(), ShouldResemble, 1)
			So(strings.Join(m[1].Namespace(), "/"), ShouldResemble, "foo/baz")
			So(m[1].Data(), ShouldResemble, "2")
		})

	})

	Convey("marshall using snap.gob", t, func() {
		m := []PluginMetricType{
			*NewPluginMetricType([]string{"foo", "bar"}, time.Now(), "", nil, nil, 1),
			*NewPluginMetricType([]string{"foo", "baz"}, time.Now(), "", nil, nil, "2"),
		}
		a, c, e := MarshalPluginMetricTypes("snap.gob", m)
		So(e, ShouldBeNil)
		So(a, ShouldNotBeNil)
		So(len(a), ShouldBeGreaterThan, 0)
		So(c, ShouldEqual, "snap.gob")

		Convey("unmarshal snap.gob", func() {
			m, e = UnmarshallPluginMetricTypes("snap.gob", a)
			So(e, ShouldBeNil)
			So(strings.Join(m[0].Namespace(), "/"), ShouldResemble, "foo/bar")
			So(m[0].Data(), ShouldResemble, 1)
			So(strings.Join(m[1].Namespace(), "/"), ShouldResemble, "foo/baz")
			So(m[1].Data(), ShouldResemble, "2")
		})

		Convey("error on bad corrupt data", func() {
			a = []byte{1, 0, 1, 1, 1, 1, 1, 0, 0, 1}
			m, e = UnmarshallPluginMetricTypes("snap.gob", a)
			So(e, ShouldNotBeNil)
			So(e.Error(), ShouldResemble, "gob: decoding into local type *[]plugin.PluginMetricType, received remote type unknown type")
		})
	})

	Convey("marshall using snap.json", t, func() {
		m := []PluginMetricType{
			*NewPluginMetricType([]string{"foo", "bar"}, time.Now(), "", nil, nil, 1),
			*NewPluginMetricType([]string{"foo", "baz"}, time.Now(), "", nil, nil, "2"),
		}
		a, c, e := MarshalPluginMetricTypes("snap.json", m)
		So(e, ShouldBeNil)
		So(a, ShouldNotBeNil)
		So(len(a), ShouldBeGreaterThan, 0)
		So(c, ShouldEqual, "snap.json")

		Convey("unmarshal snap.json", func() {
			m, e = UnmarshallPluginMetricTypes("snap.json", a)
			So(e, ShouldBeNil)
			So(strings.Join(m[0].Namespace(), "/"), ShouldResemble, "foo/bar")
			So(m[0].Data(), ShouldResemble, float64(1))
			So(strings.Join(m[1].Namespace(), "/"), ShouldResemble, "foo/baz")
			So(m[1].Data(), ShouldResemble, "2")
		})

		Convey("error on bad corrupt data", func() {
			a = []byte{1, 0, 1, 1, 1, 1, 1, 0, 0, 1}
			m, e = UnmarshallPluginMetricTypes("snap.json", a)
			So(e, ShouldNotBeNil)
			So(e.Error(), ShouldResemble, "invalid character '\\x01' looking for beginning of value")
		})
	})

	Convey("error on unmarshall using bad content type", t, func() {
		m := []PluginMetricType{
			*NewPluginMetricType([]string{"foo", "bar"}, time.Now(), "", nil, nil, 1),
			*NewPluginMetricType([]string{"foo", "baz"}, time.Now(), "", nil, nil, "2"),
		}
		a, c, e := MarshalPluginMetricTypes("snap.json", m)
		So(e, ShouldBeNil)
		So(a, ShouldNotBeNil)
		So(len(a), ShouldBeGreaterThan, 0)
		So(c, ShouldEqual, "snap.json")

		m, e = UnmarshallPluginMetricTypes("snap.wat", a)
		So(e, ShouldNotBeNil)
		So(e.Error(), ShouldEqual, "invalid snap content type for unmarshalling: snap.wat")
		So(m, ShouldBeNil)
	})

	Convey("swap from snap.gob to snap.json", t, func() {
		m := []PluginMetricType{
			*NewPluginMetricType([]string{"foo", "bar"}, time.Now(), "", nil, nil, 1),
			*NewPluginMetricType([]string{"foo", "baz"}, time.Now(), "", nil, nil, "2"),
		}
		a, c, e := MarshalPluginMetricTypes("snap.gob", m)
		So(e, ShouldBeNil)
		So(a, ShouldNotBeNil)
		So(len(a), ShouldBeGreaterThan, 0)
		So(c, ShouldEqual, "snap.gob")

		b, c, e := SwapPluginMetricContentType(c, "snap.json", a)
		So(e, ShouldBeNil)
		So(c, ShouldResemble, "snap.json")
		So(b, ShouldNotBeNil)

		m, e = UnmarshallPluginMetricTypes(c, b)
		So(e, ShouldBeNil)
		So(strings.Join(m[0].Namespace(), "/"), ShouldResemble, "foo/bar")
		So(m[0].Data(), ShouldResemble, float64(1))
		So(strings.Join(m[1].Namespace(), "/"), ShouldResemble, "foo/baz")
		So(m[1].Data(), ShouldResemble, "2")
	})

	Convey("swap from snap.json to snap.*", t, func() {
		m := []PluginMetricType{
			*NewPluginMetricType([]string{"foo", "bar"}, time.Now(), "", nil, nil, 1),
			*NewPluginMetricType([]string{"foo", "baz"}, time.Now(), "", nil, nil, "2"),
		}
		a, c, e := MarshalPluginMetricTypes("snap.json", m)
		So(e, ShouldBeNil)
		So(a, ShouldNotBeNil)
		So(len(a), ShouldBeGreaterThan, 0)
		So(c, ShouldEqual, "snap.json")

		b, c, e := SwapPluginMetricContentType(c, "snap.*", a)
		So(e, ShouldBeNil)
		So(c, ShouldResemble, "snap.gob")
		So(b, ShouldNotBeNil)

		m, e = UnmarshallPluginMetricTypes(c, b)
		So(e, ShouldBeNil)
		So(strings.Join(m[0].Namespace(), "/"), ShouldResemble, "foo/bar")
		So(m[0].Data(), ShouldResemble, float64(1))
		So(strings.Join(m[1].Namespace(), "/"), ShouldResemble, "foo/baz")
		So(m[1].Data(), ShouldResemble, "2")
	})

	Convey("swap from snap.json to snap.gob", t, func() {
		m := []PluginMetricType{
			*NewPluginMetricType([]string{"foo", "bar"}, time.Now(), "", nil, nil, 1),
			*NewPluginMetricType([]string{"foo", "baz"}, time.Now(), "", nil, nil, "2"),
		}
		a, c, e := MarshalPluginMetricTypes("snap.json", m)
		So(e, ShouldBeNil)
		So(a, ShouldNotBeNil)
		So(len(a), ShouldBeGreaterThan, 0)
		So(c, ShouldEqual, "snap.json")

		b, c, e := SwapPluginMetricContentType(c, "snap.gob", a)
		So(e, ShouldBeNil)
		So(c, ShouldResemble, "snap.gob")
		So(b, ShouldNotBeNil)

		m, e = UnmarshallPluginMetricTypes(c, b)
		So(e, ShouldBeNil)
		So(strings.Join(m[0].Namespace(), "/"), ShouldResemble, "foo/bar")
		So(m[0].Data(), ShouldResemble, float64(1))
		So(strings.Join(m[1].Namespace(), "/"), ShouldResemble, "foo/baz")
		So(m[1].Data(), ShouldResemble, "2")
	})

	Convey("error on bad content type to swap", t, func() {
		m := []PluginMetricType{
			*NewPluginMetricType([]string{"foo", "bar"}, time.Now(), "", nil, nil, 1),
			*NewPluginMetricType([]string{"foo", "baz"}, time.Now(), "", nil, nil, "2"),
		}
		a, c, e := MarshalPluginMetricTypes("snap.json", m)
		So(e, ShouldBeNil)
		So(a, ShouldNotBeNil)
		So(len(a), ShouldBeGreaterThan, 0)
		So(c, ShouldEqual, "snap.json")

		b, c, e := SwapPluginMetricContentType("snap.wat", "snap.gob", a)
		So(e, ShouldNotBeNil)
		So(e.Error(), ShouldResemble, "invalid snap content type for unmarshalling: snap.wat")
		So(b, ShouldBeNil)
	})
}
