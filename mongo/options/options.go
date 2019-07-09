package options

import "go.mongodb.org/mongo-driver/mongo/options"

func Delete() *options.DeleteOptions {
	return options.Delete()
}

func FindOne() *options.FindOneOptions {
	return options.FindOne()
}

func Find() *options.FindOptions {
	return options.Find()
}

func Client() *options.ClientOptions {
	return options.Client()
}

func Replace() *options.ReplaceOptions {
	return options.Replace()
}

func InsertOne() *options.InsertOneOptions {
	return options.InsertOne()
}

func InsertMany() *options.InsertManyOptions {
	return options.InsertMany()
}

func Transaction() *options.TransactionOptions {
	return options.Transaction()
}

func Collection() *options.CollectionOptions {
	return options.Collection()
}

func Database() *options.DatabaseOptions {
	return options.Database()
}

func Index() *options.IndexOptions {
	return options.Index()
}

func Count() *options.CountOptions {
	return options.Count()
}

func Update() *options.UpdateOptions {
	return options.Update()
}

func Aggregate() *options.AggregateOptions {
	return options.Aggregate()
}

func BulkWrite() *options.BulkWriteOptions {
	return options.BulkWrite()
}

func ChangeStream() *options.ChangeStreamOptions {
	return options.ChangeStream()
}

func CreateIndexes() *options.CreateIndexesOptions {
	return options.CreateIndexes()
}

func Distinct() *options.DistinctOptions {
	return options.Distinct()
}

func DropIndexes() *options.DropIndexesOptions {
	return options.DropIndexes()
}

func EstimatedDocumentCount() *options.EstimatedDocumentCountOptions {
	return options.EstimatedDocumentCount()
}

func Session() *options.SessionOptions {
	return options.Session()
}

func RunCmd() *options.RunCmdOptions {
	return options.RunCmd()
}

func ListCollections() *options.ListCollectionsOptions {
	return options.ListCollections()
}

func ListDatabases() *options.ListDatabasesOptions {
	return options.ListDatabases()
}

func ListIndexes() *options.ListIndexesOptions {
	return options.ListIndexes()
}

func GridFSBucket() *options.BucketOptions {
	return options.GridFSBucket()
}

func GridFSName() *options.NameOptions {
	return options.GridFSName()
}

func FindOneAndDelete() *options.FindOneAndDeleteOptions {
	return options.FindOneAndDelete()
}

func FindOneAndReplace() *options.FindOneAndReplaceOptions {
	return options.FindOneAndReplace()
}

func FindOneAndUpdate() *options.FindOneAndUpdateOptions {
	return options.FindOneAndUpdate()
}

func MergeGridFSFindOptions() *options.GridFSFindOptions {
	return options.GridFSFind()
}

func GridFSUpload() *options.UploadOptions {
	return options.GridFSUpload()
}
