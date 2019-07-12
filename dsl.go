package go_cypherdsl

import (
	"errors"
	"fmt"
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
)

type ConnectionConfig struct {
	Username string
	Password string
	Host string
	Port int
	PoolSize int
}

func (c *ConnectionConfig) ConnectionString() string{
	return fmt.Sprintf("bolt://%s:%s@%s:%v", c.Username, c.Password, c.Host, c.Port)
}

var connPool bolt.DriverPool

var isInitialized = false

func Init(connection *ConnectionConfig) error{
	if isInitialized{
		return errors.New("already initialized")
	}

	if connection == nil{
		return errors.New("connection can not be nil")
	}

	//if pool size isn't set
	if connection.PoolSize <= 0{
		//set default to 15
		connection.PoolSize = 15
	}

	var err error
	connPool, err = bolt.NewDriverPool(connection.ConnectionString(), connection.PoolSize)
	if err != nil{
		return err
	}

	isInitialized = true

	return nil
}

type Session struct {
	conn bolt.Conn
	tx bolt.Tx
}

func NewSession() *Session{
	return new(Session)
}

func (s *Session) Begin() error{
	if !isInitialized{
		return errors.New("cypher dsl not initialized")
	}

	var err error
	if s.conn == nil{
		s.conn, err = connPool.OpenPool()
		if err != nil{
			return err
		}
	}

	s.tx, err = s.conn.Begin()
	if err != nil{
		return err
	}

	return nil
}

func (s *Session) Rollback() error{
	if !isInitialized{
		return errors.New("cypher dsl not initialized")
	}

	if s.tx == nil{
		return errors.New("transaction not initialized")
	}

	err := s.tx.Rollback()
	if err != nil{
		return err
	}

	//set transaction to nil for other logic
	s.tx = nil

	return nil
}

func (s *Session) Commit() error{
	if !isInitialized{
		return errors.New("cypher dsl not initialized")
	}

	if s.tx == nil{
		return errors.New("transaction not initialized")
	}

	err := s.tx.Commit()
	if err != nil{
		return err
	}

	//set transaction to nil for other logic
	s.tx = nil

	return nil
}

func (s *Session) Close() error{
	if !isInitialized{
		return errors.New("cypher dsl not initialized")
	}

	if s.conn == nil{
		return errors.New("connection not open")
	}

	return s.conn.Close()
}

//to do a query
func (s *Session) Query() Cypher{
	if !isInitialized{
		return &QueryBuilder{
			errors: []error{errors.New("cypher dsl not initialized")},
		}
	}

	var err error

	//if the connection is not initialized, initialize it
	if s.conn == nil{
		s.conn, err = connPool.OpenPool()
		if err != nil{
			return &QueryBuilder{
				errors: []error{err},
			}
		}
	}

	//if the transaction is nil, tell the query builder to make its own connection and transaction
	if s.tx == nil{
		return &QueryBuilder{
			conn: nil,
		}
	}

	return&QueryBuilder{
		conn: s.conn,
	}
}