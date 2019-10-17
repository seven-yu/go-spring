/*
 * Copyright 2012-2019 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package SpringCore_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/go-spring/go-spring/spring-core"
	"github.com/go-spring/go-spring/spring-core/testdata/bar"
	otherfoo "github.com/go-spring/go-spring/spring-core/testdata/bar/foo"
	"github.com/go-spring/go-spring/spring-core/testdata/foo"
	"github.com/go-spring/go-spring/spring-utils"
	"github.com/stretchr/testify/assert"
)

type EvenetHandler interface {
	Notify() string
}

type FooEventHandler struct {
	Name    string
	*Config `autowire:""`
}

func (e *FooEventHandler) Notify() string {
	return fmt.Sprintf("I am %s , config port: %d !", e.Name, e.Config.Port)
}

type BarEventHandler struct {
	Name    string
	*Config `autowire:""`
}

func (e *BarEventHandler) Notify() string {
	return fmt.Sprintf("I am %s , config port: %d !", e.Name, e.Config.Port)
}

type EventDemo struct {
	Handlers []EvenetHandler `autowire:""`
}

type Config struct {
	Port int32 `value:"${test.spring.core.config.port:=1234}"`
}

func TestInterfaceSliceBeanWire(t *testing.T) {
	ctx := SpringCore.NewDefaultSpringContext()

	c := &Config{Port: 9090}
	ctx.RegisterSingletonBean(c)

	f := &FooEventHandler{Name: "Foo"}
	ctx.RegisterSingletonBean(f)

	b := &BarEventHandler{Name: "Bar"}
	ctx.RegisterSingletonBean(b)

	demo := new(EventDemo)
	ctx.RegisterSingletonBean(demo)

	_ = ctx.AutoWireBeans()

	eventDemo := ctx.FindBeanByType(&EventDemo{}).(*EventDemo)

	assert.Equal(t, len(eventDemo.Handlers), 2)

	for _, h := range eventDemo.Handlers {
		t.Log(h.Notify())
	}
}

type ConfigsDemo struct {
	Configs []*Config `autowire:""`
	Nums    []int     `autowire:""`
	Strings []string  `autowire:""`
}

func TestStructSliceWire(t *testing.T) {
	ctx := SpringCore.NewDefaultSpringContext()

	cfg := new(Config)
	ctx.RegisterSingletonBean(cfg)

	nums := []int{1, 2}
	ctx.RegisterSingletonBean(nums)

	strs := []string{"a", "b"}
	ctx.RegisterSingletonBean(strs)

	cfgs := []*Config{
		{
			Port: 8080,
		},
		{
			Port: 8081,
		},
	}
	ctx.RegisterSingletonBean(cfgs)

	demo := new(ConfigsDemo)
	ctx.RegisterSingletonBean(demo)

	_ = ctx.AutoWireBeans()

	cfgsDemo := ctx.FindBeanByType(&ConfigsDemo{}).(*ConfigsDemo)
	assert.Equal(t, cfgsDemo.Configs[0].Port, int32(8080))
	assert.Equal(t, cfgsDemo.Nums[0], 1)
	assert.Equal(t, cfgsDemo.Strings[0], "a")
	assert.Equal(t, cfgsDemo.Configs[1].Port, int32(8081))
	assert.Equal(t, cfgsDemo.Nums[1], 2)
	assert.Equal(t, cfgsDemo.Strings[1], "b")

}

func TestStructSlice(t *testing.T) {

	cfgs := []*Config{
		{
			Port: 8080,
		},
		{
			Port: 8081,
		},
	}

	typeOf := reflect.TypeOf(&cfgs)

	beanName := fmt.Sprintf(
		"%s.%s",
		strings.Replace(typeOf.Elem().PkgPath(), "/", ".", -1),
		typeOf.Elem().Name(),
	)

	t.Log(beanName)

}

func TestValueWire(t *testing.T) {

	type People struct {
		FirstName string `value:"${people.first_name}"`
		LastName  string `value:"${people.last_name:=Green}"`

		Uint  uint64  `value:"${people.uint:=1}"`
		Int   int64   `value:"${people.int:=-1}"`
		Float float64 `value:"${people.float:=3.14}"`
		Bool  bool    `value:"${people.bool:=true}"`
	}

	ctx := SpringCore.NewDefaultSpringContext()

	p := new(People)
	ctx.RegisterSingletonBean(p)

	ctx.SetProperties("people.first_name", "Jim")

	if err := ctx.AutoWireBeans(); err != nil {
		panic(err)
	}

	var find *People
	ctx.GetBeanByType(&find)

	assert.Equal(t, "Jim", find.FirstName)
	assert.Equal(t, "Green", find.LastName)
	assert.Equal(t, 3.14, find.Float)
	assert.Equal(t, int64(-1), find.Int)
	assert.Equal(t, uint64(1), find.Uint)
	assert.Equal(t, true, find.Bool)

	fmt.Println(SpringUtils.ToJson(p))
}

func TestBeanWire(t *testing.T) {

	type Config struct {
		Name string
	}

	type DataSource struct {
		Url string
	}

	type Application struct {
		Config     *Config     `autowire:""`
		DataSource *DataSource `autowire:"ds"`
	}

	ctx := SpringCore.NewDefaultSpringContext()

	app := new(Application)
	ctx.RegisterSingletonBean(app)

	cfg := &Config{Name: "application.cfg"}
	ctx.RegisterSingletonBean(cfg)

	ds := &DataSource{
		Url: "mysql:127.0.0.1...",
	}

	ctx.RegisterSingletonBean(ds)
	ctx.RegisterSingletonNameBean("ds", ds)

	barBean := new(foo.Demo)
	fooBean := new(bar.Demo)
	ctx.RegisterSingletonBean(barBean)
	ctx.RegisterSingletonBean(fooBean)

	if e := ctx.AutoWireBeans(); e != nil {
		t.Error(e)
	}

	for _, v := range ctx.GetAllBeanNames() {
		t.Logf("bean name : %v", v)
	}

	var (
		f foo.Demo
		b bar.Demo
	)

	foundFooBean := ctx.FindBeanByType(&f)
	foundFarBean := ctx.FindBeanByType(&b)

	assert.NotEqual(t, foundFarBean, foundFooBean)

	fmt.Println(SpringUtils.ToJson(app))
}

func TestSlice(t *testing.T) {
	f := &FooEventHandler{Name: "foo"}
	b := &BarEventHandler{Name: "bar"}
	var handlers []EvenetHandler
	slice := reflect.MakeSlice(reflect.TypeOf(handlers), 0, 0)
	slice = reflect.Append(slice, reflect.ValueOf(f), reflect.ValueOf(b))
	t.Log(slice)
}

func TestField(t *testing.T) {
	f := &FooEventHandler{Name: "foo"}
	b := &BarEventHandler{Name: "bar"}
	var handlers []EvenetHandler
	handlers = append(handlers, f, b)
	demo := &EventDemo{Handlers: handlers}

	demoType := reflect.TypeOf(demo)

	filed := demoType.Elem().Field(0)

	filedInterface := filed.Type.Elem()

	t.Log(filed.Type.Name())

	slice := reflect.MakeSlice(filed.Type, 2, 2)
	reflect.Append(slice, reflect.ValueOf(f))
	t.Log(slice)

	t.Log(reflect.TypeOf(f).Implements(filedInterface))
	t.Log(filedInterface)

}

func TestRegiserSingletonBean(t *testing.T) {
	assert.Panics(t, func() {
		ctx := SpringCore.NewDefaultSpringContext()

		cfg := new(Config)
		ctx.RegisterSingletonBean(cfg)

		cfg2 := &Config{Port: int32(1111)}
		ctx.RegisterSingletonBean(cfg2)
	})
}

func TestFindSliceBean(t *testing.T) {
	ctx := SpringCore.NewDefaultSpringContext()

	ctx.RegisterSingletonBean(&Config{Port: 1234})

	cfgs := []*Config{
		{
			Port: 8080,
		},
		{
			Port: 8081,
		},
	}
	ctx.RegisterSingletonBean(cfgs)

	var find []*Config
	ctx.GetBeanByType(&find)

	assert.Equal(t, find[0].Port, int32(8080))
	assert.Equal(t, find[1].Port, int32(8081))

	var cfg *Config
	ctx.GetBeanByType(&cfg)
	assert.Equal(t, cfg.Port, int32(1234))
}

func TestGetBeanUname(t *testing.T) {
	beanUname := SpringCore.GetTypeName(reflect.TypeOf(&foo.Demo{}))
	otherPkgBeanUname := SpringCore.GetTypeName(reflect.TypeOf(&bar.Demo{}))
	assert.NotEqual(t, beanUname, otherPkgBeanUname)
}

func TestGetBeanName(t *testing.T) {
	ctx := SpringCore.NewDefaultSpringContext()
	ctx.RegisterSingletonBean(new(Config))
	var cfg *Config
	ctx.FindBeanByType(cfg)
	t.Log(SpringCore.GetBeanName(cfg))

	ctx.GetAllBeanNames()
}

func TestRegisterSingletonBean(t *testing.T) {
	ctx := SpringCore.NewDefaultSpringContext()

	f := &foo.Demo{}
	b := &bar.Demo{}
	otherf := &otherfoo.Demo{}

	ctx.RegisterSingletonBean(f)
	ctx.RegisterSingletonBean(b)
	ctx.RegisterSingletonBean(otherf)

	_ = ctx.AutoWireBeans()

	t.Log(ctx.GetAllBeanNames())

}

func TestGetTypeName(t *testing.T) {

	ctx := SpringCore.NewDefaultSpringContext()

	t.Run("struct ptr slice bean", func(t *testing.T) {
		cfgs := []*Config{
			{
				Port: 8080,
			},
			{
				Port: 8081,
			},
		}
		ctx.RegisterSingletonBean(cfgs)
	})

	t.Run("struct slice", func(t *testing.T) {
		cfgs2 := []Config{
			{
				Port: 8080,
			},
			{
				Port: 8081,
			},
		}
		ctx.RegisterSingletonBean(cfgs2)
	})

	t.Run("struct ptr bean", func(t *testing.T) {
		ctx.RegisterSingletonBean(new(Config))
	})

	t.Run("built in type bean", func(t *testing.T) {
		nums := []int{1, 2}
		ctx.RegisterSingletonBean(nums)

		strs := []string{"a", "b"}
		ctx.RegisterSingletonBean(strs)
	})

	_ = ctx.AutoWireBeans()

	t.Log(ctx.GetAllBeanNames())

}

func TestPkgPath(t *testing.T) {

	num := 1
	t.Log(reflect.TypeOf(num).PkgPath())

	nums := []int{1, 2}
	t.Log(reflect.TypeOf(nums).Elem().PkgPath())
}

// TODO beanmap的key可以参考这个
type key struct {
	t reflect.Type

	// 允许重复注册同一种类型，即 => type相同，name不同
	name string
}

func TestReflectTypeAsKey(t *testing.T) {
	f := &foo.Demo{}
	b := &bar.Demo{}
	otherf := &otherfoo.Demo{}
	nums := []int{1, 2, 3}
	strs := []string{"a", "b", "c"}
	cfgs := []*Config{
		{
			Port: 8080,
		},
		{
			Port: 8081,
		},
	}
	cfgsNotPtr := []Config{
		{
			Port: 8080,
		},
		{
			Port: 8081,
		},
	}

	var beanMap = make(map[key]reflect.Value)

	assert.NotEqual(t, reflect.TypeOf(f).Elem(), reflect.TypeOf(otherf).Elem())
	assert.Equal(t, reflect.TypeOf(f).Name(), reflect.TypeOf(otherf).Name())

	beanMap[key{t: reflect.TypeOf(f)}] = reflect.ValueOf(f)
	beanMap[key{t: reflect.TypeOf(b)}] = reflect.ValueOf(b)
	beanMap[key{t: reflect.TypeOf(otherf)}] = reflect.ValueOf(otherf)
	beanMap[key{t: reflect.TypeOf(nums)}] = reflect.ValueOf(nums)
	beanMap[key{t: reflect.TypeOf(strs)}] = reflect.ValueOf(strs)
	beanMap[key{t: reflect.TypeOf(cfgs)}] = reflect.ValueOf(cfgs)
	beanMap[key{t: reflect.TypeOf(cfgsNotPtr)}] = reflect.ValueOf(cfgsNotPtr)

	assert.Equal(t, 7, len(beanMap))
}
