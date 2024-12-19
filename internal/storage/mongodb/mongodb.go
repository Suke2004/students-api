package mongodb

import (
	"context"
	"time"

	"github.com/Suke2004/students-api/internal/config"
	"github.com/Suke2004/students-api/internal/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func New(cfg *config.Config) (*MongoDB, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(cfg.MongoDB.URI))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	collection := client.Database(cfg.MongoDB.Database).Collection("students")
	return &MongoDB{
		client:     client,
		collection: collection,
	}, nil
}


// CreateStudent inserts a new student document.
func (m *MongoDB) CreateStudent(name string, email string, age int) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	student := types.Student{
		Name:  name,
		Email: email,
		Age:   age,
	}
	result, err := m.collection.InsertOne(ctx, student)
	if err != nil {
		return 0, err
	}

	// Return the inserted ID as an int64 (not typical in MongoDB)
	id := result.InsertedID.(primitive.ObjectID).Timestamp().Unix()
	return id, nil
}

// GetStudentById retrieves a student document by ID.
func (m *MongoDB) GetStudentById(id int64) (types.Student, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Convert ID to ObjectID
	objectID := primitive.NewObjectIDFromTimestamp(time.Unix(id, 0))

	var student types.Student
	err := m.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&student)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return types.Student{}, nil
		}
		return types.Student{}, err
	}

	return student, nil
}

// GetStudent retrieves all student documents.
func (m *MongoDB) GetStudent() ([]types.Student, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := m.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var students []types.Student
	for cursor.Next(ctx) {
		var student types.Student
		err := cursor.Decode(&student)
		if err != nil {
			return nil, err
		}
		students = append(students, student)
	}

	return students, nil
}

// DeleteStudentById deletes a student document by ID.
func (m *MongoDB) DeleteStudentById(id int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Convert ID to ObjectID
	objectID := primitive.NewObjectIDFromTimestamp(time.Unix(id, 0))

	_, err := m.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

// DeleteAllStudents deletes all student documents.
func (m *MongoDB) DeleteAllStudents() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := m.collection.DeleteMany(ctx, bson.M{})
	return err
}

// UpdateStudentById updates a student document by ID.
func (m *MongoDB) UpdateStudentById(id int64, name, email string, age int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Convert ID to ObjectID
	objectID := primitive.NewObjectIDFromTimestamp(time.Unix(id, 0))

	// Update document fields
	update := bson.M{
		"$set": bson.M{
			"name":  name,
			"email": email,
			"age":   age,
		},
	}

	_, err := m.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}
