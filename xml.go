package main

type XMLResponse struct {
	ListBucketResult string
	Name             string
	Prefix           string
	Maker            string
	Delimiter        string
	IsTruncated      string
	CommonPrefixes   []CommonPrefix
}

type CommonPrefix struct {
	Prefix string
}
