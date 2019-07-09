package go_cypherdsl

import (
	"errors"
	"fmt"
	neo "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"strings"
)

type QueryBuilder struct {
	Start *queryPartNode
	Current *queryPartNode
	errors []error

	conn neo.Conn
}

func QB() *QueryBuilder{
	return &QueryBuilder{}
}

func (q *QueryBuilder) addNext(s string) {
	node := &queryPartNode{
		Part: s,
	}

	if q.Start == nil{
		q.Start = node
		q.Current = node
	} else {
		q.Current.Next = node
		q.Current = node
	}
}

func (q *QueryBuilder) addError(err error){
	if q.errors == nil{
		q.errors = []error{}
	}

	q.errors = append(q.errors, err)
}

func (q *QueryBuilder) hasErrors() bool{
	return q.errors != nil && len(q.errors) > 0
}

type queryPartNode struct {
	Part string
	Next *queryPartNode
}

func (q *QueryBuilder) Match(p *PathBuilder) Cypher {
	if p == nil{
		q.addError(errors.New("path can not be nil"))
		return q
	}

	query, err := p.ToCypher()
	if err != nil{
		q.addError(err)
		return q
	}

	q.addNext("MATCH " + query)
	return q
}

func (q *QueryBuilder) Create(c CreateQuery, err error) Cypher {
	if err != nil{
		q.addError(err)
		return q
	}

	q.addNext("CREATE " + string(c))
	return q
}

func (q *QueryBuilder) Where(cb ConditionOperator) Cypher {
	if cb == nil{
		q.addError(errors.New("condition builder can not be nil"))
		return q
	}

	w, err := cb.Build()
	if err != nil{
		q.addError(err)
		return q
	}

	q.addNext("WHERE " + string(w))
	return q
}

func (q *QueryBuilder) Merge(mergeConf *MergeConfig) Cypher {
	if mergeConf == nil{
		q.addError(errors.New("mergeConf can not be nil"))
		return q
	}
	cypher, err := mergeConf.ToString()
	if err != nil{
		q.addError(err)
		return q
	}

	q.addNext("MERGE " + cypher)

	return q
}

func (q *QueryBuilder) Return(parts ...ReturnPart) Cypher {
	str, err := NewReturnClause(parts...)
	if err != nil{
		q.addError(err)
		return q
	}

	q.addNext(string(str))
	return q
}

func (q *QueryBuilder) Delete(detach bool, params ...string) Cypher {
	cypher, err := deleteToString(detach, params...)
	if err != nil{
		q.addError(err)
		return q
	}

	q.addNext(cypher)
	return q
}

func (q *QueryBuilder) Set(sets ...SetConfig) Cypher {
	if len(sets) == 0{
		q.addError(errors.New("sets can not be empty"))
		return q
	}

	query := "SET"

	for _, setStmt := range sets{
		str, err := setStmt.ToString()
		if err != nil{
			q.addError(err)
			return q
		}

		query += fmt.Sprintf(" %s,", str)
	}

	q.addNext(strings.TrimSuffix(query, ","))
	return q
}

func (q *QueryBuilder) Remove(removes ...RemoveConfig) Cypher {
	if len(removes) == 0{
		q.addError(errors.New("removes can not be empty"))
	}

	query := "REMOVE"

	for _, remove := range removes{
		str, err := remove.ToString()
		if err != nil{
			q.addError(err)
			return q
		}
		query += fmt.Sprintf(" %s,", str)
	}

	q.addNext(strings.TrimSuffix(query, ","))
	return q
}

func (q *QueryBuilder) OrderBy(orderBys ...OrderByConfig) Cypher{
	if len(orderBys) == 0{
		q.addError(errors.New("removes can not be empty"))
	}

	query := "ORDER BY"

	for _, orders := range orderBys{
		str, err := orders.ToString()
		if err != nil{
			q.addError(err)
			return q
		}
		query += fmt.Sprintf(" %s,", str)
	}

	q.addNext(strings.TrimSuffix(query, ","))
	return q
}

func (q *QueryBuilder) Limit(num int) Cypher{
	q.addNext(fmt.Sprintf("LIMIT %v", num))
	return q
}

func (q *QueryBuilder) Query(params map[string]interface{}) (neo.Rows, error) {
	query, err := q.build()
	if err != nil{
		return nil, err
	}

	return q.conn.QueryNeo(query, params)
}

func (q *QueryBuilder) QueryStruct(params map[string]interface{}, respObj interface{}) (neo.Rows, error) {
	query, err := q.build()
	if err != nil{
		return nil, err
	}

	res, err := q.conn.QueryNeo(query, params)
	if err != nil{
		return nil, err
	}

	//todo actually handle

	return res, nil
}

func (q *QueryBuilder) Exec(params map[string]interface{}) (neo.Result, error){
	query, err := q.build()
	if err != nil{
		return nil, err
	}

	return q.conn.ExecNeo(query, params)
}

func (q *QueryBuilder) ToCypher() (string, error){
	return q.build()
}

func (q *QueryBuilder) build() (string, error){
	//fail if errors are found
	if q.hasErrors(){
		str := "errors found: "
		for _, err := range q.errors{
			str += err.Error() + ";"
		}

		str = strings.TrimSuffix(str, ";") + fmt.Sprintf(" -- total errors (%v)", len(q.errors))
		return "", errors.New(str)
	}

	if q.Start == nil || q.Current == nil{
		return "", errors.New("no nodes were added")
	}

	query := ""

	cur := q.Start

	for {
		if cur == nil{
			break
		}

		query += cur.Part + " "

		cur = cur.Next
	}

	return strings.TrimSuffix(query, " "), nil
}