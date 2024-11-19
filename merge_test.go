package go_cypherdsl

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMergeSetConfig_ToString(t *testing.T) {
	var err error
	var cypher string

	req := require.New(t)
	t1 := MergeSetConfig{Name: "test", Member: "ttt", Target: 1}
	t2 := MergeSetConfig{Name: "test", Member: "ttt", TargetFunction: &FunctionConfig{
		Name: "test",
	}}

	t3 := MergeSetConfig{Name: "test"}
	t4 := MergeSetConfig{}
	t5 := MergeSetConfig{Name: "test", Member: "ttt"}
	t6 := MergeSetConfig{Name: "test", Member: "ttt", TargetFunction: &FunctionConfig{Name: "test"}, Target: 1}

	params, err := ParamsFromMap(map[string]interface{}{
		"name": "li",
		"age":  14,
	})
	req.Nil(err)
	t7 := MergeSetConfig{Name: "n", TargetMap: params, Operation: SetMutate}

	//name member normal target
	cypher, err = t1.ToString()
	req.Nil(err)
	req.EqualValues("test.ttt = 1", cypher)

	//name member target function
	cypher, err = t2.ToString()
	req.Nil(err)
	req.EqualValues("test.ttt = test()", cypher)

	//error - member not defined
	_, err = t3.ToString()
	req.NotNil(err)

	//error - member and name not defined
	_, err = t4.ToString()
	req.NotNil(err)

	//error - target and target function not defined
	_, err = t5.ToString()
	req.NotNil(err)

	//error - target and target function defined
	_, err = t6.ToString()
	req.NotNil(err)

	cypher, err = t7.ToString()
	req.Nil(err)
	req.EqualValues("n += {name:'li',age:14}", cypher)
}

func TestMergeConfig_ToString(t *testing.T) {
	t1 := MergeConfig{Path: "test"}

	t2 := MergeConfig{Path: "test", OnCreate: &MergeSetConfig{
		Name:   "test",
		Member: "ttt",
		Target: 1,
	}}

	t3 := MergeConfig{Path: "test", OnMatch: &MergeSetConfig{
		Name:   "test",
		Member: "ttt",
		Target: 1,
	}}

	t4 := MergeConfig{Path: "test", OnMatch: &MergeSetConfig{
		Name:   "test",
		Member: "ttt",
		Target: 1,
	}, OnCreate: &MergeSetConfig{
		Name:   "test",
		Member: "tt1",
		Target: 2,
	}}

	t5 := MergeConfig{}

	req := require.New(t)
	var err error
	var cypher string

	//only merge
	cypher, err = t1.ToString()
	req.Nil(err)
	req.EqualValues("test", cypher)

	//merge with on create
	cypher, err = t2.ToString()
	req.Nil(err)
	req.EqualValues("test ON CREATE SET test.ttt = 1", cypher)

	//merge with on match
	cypher, err = t3.ToString()
	req.Nil(err)
	req.EqualValues("test ON MATCH SET test.ttt = 1", cypher)

	//merge with on create and on match
	cypher, err = t4.ToString()
	req.Nil(err)
	req.EqualValues("test ON CREATE SET test.tt1 = 2 ON MATCH SET test.ttt = 1", cypher)

	//error - path not defined
	_, err = t5.ToString()
	req.NotNil(err)
}
