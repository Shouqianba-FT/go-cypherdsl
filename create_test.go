package go_cypherdsl

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewConstraint(t *testing.T) {
	req := require.New(t)

	//test nil
	_, err := NewConstraint(nil)
	req.NotNil(err)

	//test empty name, type fields
	errConf := &ConstraintConfig{Name: "", Type: "s", Field: "s"}

	//name empty
	_, err = NewConstraint(errConf)
	req.NotNil(err)

	//name and type
	errConf.Type = ""
	_, err = NewConstraint(errConf)
	req.NotNil(err)

	//name type and fields
	errConf.Field = ""
	_, err = NewConstraint(errConf)
	req.NotNil(err)

	//test unique and exists both false, both true
	errConf.Unique = false
	errConf.Exists = false
	_, err = NewConstraint(errConf)
	req.NotNil(err)

	errConf.Unique = true
	errConf.Exists = true
	_, err = NewConstraint(errConf)
	req.NotNil(err)

	//test unique
	testConf := errConf

	testConf.Exists = false
	testConf.Name = "t"
	testConf.Type = "s"
	testConf.Field = "s"
	cypher, err := NewConstraint(testConf)
	req.Nil(err)
	req.EqualValues("CONSTRAINT IF NOT EXISTS FOR (t:s) REQUIRE t.s IS UNIQUE", cypher.ToString())

	//test exists
	testConf.Exists = true
	testConf.Unique = false
	cypher, err = NewConstraint(testConf)
	req.Nil(err)
	req.EqualValues("CONSTRAINT IF NOT EXISTS FOR (t:s) REQUIRE t.s IS NOT NULL", cypher.ToString())
}

func TestNewIndex(t *testing.T) {
	req := require.New(t)

	//check nil
	_, err := NewIndex(nil)
	req.NotNil(err)

	conf := &IndexConfig{Fields: nil}

	//check empty type
	_, err = NewIndex(conf)
	req.NotNil(err)

	//check nil fields
	conf.PathBuilder = Path().V(V{Name: "n", Type: "test"}).Build()
	conf.Index = "i_test_id"
	conf.Name = "n"
	_, err = NewIndex(conf)
	req.NotNil(err)

	//check empty fields
	conf.Fields = []string{}
	_, err = NewIndex(conf)
	req.NotNil(err)

	//check single index
	conf.Fields = []string{"one"}
	cypher, err := NewIndex(conf)
	req.Nil(err)
	req.EqualValues("INDEX i_test_id IF NOT EXISTS FOR (n:test) ON (n.one)", cypher.ToString())

	//check composite index
	conf.Fields = []string{"one", "two"}
	cypher, err = NewIndex(conf)
	req.Nil(err)
	req.EqualValues("INDEX i_test_id IF NOT EXISTS FOR (n:test) ON (n.one,n.two)", cypher.ToString())
}

func TestNewNode(t *testing.T) {
	req := require.New(t)

	_, err := NewNode(nil)
	req.NotNil(err)

	query, err := NewNode(Path().V(V{}).Build())
	req.Nil(err)
	req.NotEqual("test", query)
}
