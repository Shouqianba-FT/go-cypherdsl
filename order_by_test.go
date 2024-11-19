package go_cypherdsl

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestOrderByConfig_ToString(t *testing.T) {
	req := require.New(t)
	var err error
	var cypher string

	//name not defined
	t1 := OrderByConfig{
		Type: "da",
	}
	_, err = t1.ToString()
	req.NotNil(err)

	//Type not defined
	t2 := OrderByConfig{
		Name: "da",
	}
	_, err = t2.ToString()
	req.Nil(err)

	//both type and name not defined
	t3 := OrderByConfig{}
	_, err = t3.ToString()
	req.NotNil(err)

	//proper
	t4 := OrderByConfig{
		Type: "n",
		Name: "name",
	}
	cypher, err = t4.ToString()
	req.Nil(err)
	req.EqualValues("n.name", cypher)

	//proper
	t5 := OrderByConfig{
		Type: "n",
		Name: "name",
		Desc: true,
	}
	cypher, err = t5.ToString()
	req.Nil(err)
	req.EqualValues("n.name DESC", cypher)

	// type not defined
	t6 := OrderByConfig{
		Name: "name",
		Desc: true,
	}
	cypher, err = t6.ToString()
	req.Nil(err)
	req.EqualValues("name DESC", cypher)

	// type and desc not defined
	t7 := OrderByConfig{
		Name: "name",
	}
	cypher, err = t7.ToString()
	req.Nil(err)
	req.EqualValues("name", cypher)
}
