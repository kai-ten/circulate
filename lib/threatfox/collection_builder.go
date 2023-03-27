package main

import (
	"crypto/sha512"
	"encoding/hex"
)

// Golang defaults to utf-8 character encoding
// Hex is the common way to represent a hash byte array
func uniqueKey(input string) string {
	hash := sha512.Sum512([]byte(input))
	hash_str := hex.EncodeToString(hash[:])
	return hash_str
}

func mapToIOC(tfIOC *ThreatfoxIOC) IOC {
	key := uniqueKey(tfIOC.ID)
	ioc := IOC{
		ID:              "ioc/" + key,
		Key:             key,
		Value:           tfIOC.Value,
		Type:            tfIOC.Type,
		ThreatType:      tfIOC.ThreatType,
		ThreatTypeDesc:  tfIOC.ThreatTypeDesc,
		ConfidenceLevel: tfIOC.ConfidenceLevel,
		FirstSeen:       tfIOC.FirstSeen,
		LastSeen:        tfIOC.LastSeen,
	}
	return ioc
}

func mapToMalware(tfIOC *ThreatfoxIOC) Malware {
	key := uniqueKey(tfIOC.Malware)
	malware := Malware{
		ID:               "malware/" + key,
		Key:              key,
		Malware:          tfIOC.Malware,
		MalwarePrintable: tfIOC.MalwarePrintable,
		MalwareAlias:     tfIOC.MalwareAlias,
		MalwareMalpedia:  tfIOC.MalwareMalpedia,
		FirstSeen:        tfIOC.FirstSeen,
		LastSeen:         tfIOC.LastSeen,
		Reference:        tfIOC.Reference,
		Tags:             tfIOC.Tags,
	}
	return malware
}

func mapToReporter(tfIOC *ThreatfoxIOC) Reporter {
	key := uniqueKey(tfIOC.Reporter)
	reporter := Reporter{
		ID:       "reporter/" + key,
		Key:      key,
		Reporter: tfIOC.Reporter,
	}
	return reporter
}

func mapToThreatSource(tfIOC *ThreatfoxIOC) ThreatSource {
	source := "ThreatFox"
	hash := uniqueKey(source)

	threatSource := ThreatSource{
		ID:    "threatSource/" + hash,
		Key:   hash,
		Value: source,
	}
	return threatSource
}

func mapVerticesToEdge(from string, to string) Edge {
	return Edge{
		From: from,
		To:   to,
	}
}

func createEdges(iocVertex IOC, threatSourceVertex ThreatSource, malwareVertex Malware, reporterVertex Reporter) (iocToMalwareEdges, iocToThreatSourceEdges, threatSourceToIocEdges, malwareToIocEdges, iocToReporterEdges, malwareToReporterEdges, reporterToIocEdges, reporterToMalwareEdges, reporterToThreatSourceEdges Edge) {
	iocToReporterEdge := mapVerticesToEdge(iocVertex.ID, reporterVertex.ID)
	iocToMalwareEdge := mapVerticesToEdge(iocVertex.ID, malwareVertex.ID)
	iocToThreatSourceEdge := mapVerticesToEdge(iocVertex.ID, threatSourceVertex.ID)
	threatSourceToIocEdge := mapVerticesToEdge(threatSourceVertex.ID, iocVertex.ID)
	malwareToIocEdge := mapVerticesToEdge(malwareVertex.ID, iocVertex.ID)
	malwareToReporterEdge := mapVerticesToEdge(malwareVertex.ID, reporterVertex.ID)
	reporterToIocEdge := mapVerticesToEdge(reporterVertex.ID, iocVertex.ID)
	reporterToMalwareEdge := mapVerticesToEdge(reporterVertex.ID, malwareVertex.ID)
	reporterToThreatSourceEdge := mapVerticesToEdge(reporterVertex.ID, threatSourceVertex.ID)

	return iocToMalwareEdge, iocToThreatSourceEdge, threatSourceToIocEdge, malwareToIocEdge, iocToReporterEdge, malwareToReporterEdge, reporterToIocEdge, reporterToMalwareEdge, reporterToThreatSourceEdge
}

// Convert JSON into logical vertices and edges
func BuildCollections(threatFoxResp ThreatFoxResponse) (vertexCollection map[string]interface{}, edgeCollection map[string][]Edge) {
	var iocVertices []IOC
	var malwareVertices []Malware
	var reporterVertices []Reporter
	var threatSourceVertices []ThreatSource

	var iocToMalwareEdges []Edge
	var iocToReporterEdges []Edge
	var malwareToIocEdges []Edge
	var malwareToReporterEdges []Edge
	var reporterToIocEdges []Edge
	var reporterToMalwareEdges []Edge
	var iocToThreatSourceEdges []Edge
	var threatSourceToIocEdges []Edge
	var reporterToThreatSourceEdges []Edge

	for _, tfioc := range threatFoxResp.Result {
		// Create Vertices
		iocVertex := mapToIOC(&tfioc)
		threatSourceVertex := mapToThreatSource(&tfioc)
		malwareVertex := mapToMalware(&tfioc)
		reporterVertex := mapToReporter(&tfioc)

		// Append to Vertex Collections
		iocVertices = append(iocVertices, iocVertex)
		threatSourceVertices = append(threatSourceVertices, threatSourceVertex)
		malwareVertices = append(malwareVertices, malwareVertex)
		reporterVertices = append(reporterVertices, reporterVertex)

		// Create and Append Edges
		iocToMalwareEdge, iocToThreatSourceEdge, threatSourceToIocEdge, malwareToIocEdge, iocToReporterEdge, malwareToReporterEdge, reporterToIocEdge, reporterToMalwareEdge, reporterToThreatSourceEdge := createEdges(iocVertex, threatSourceVertex, malwareVertex, reporterVertex)
		iocToMalwareEdges = append(iocToMalwareEdges, iocToMalwareEdge)
		iocToThreatSourceEdges = append(iocToThreatSourceEdges, iocToThreatSourceEdge)
		threatSourceToIocEdges = append(threatSourceToIocEdges, threatSourceToIocEdge)
		malwareToIocEdges = append(malwareToIocEdges, malwareToIocEdge)
		iocToReporterEdges = append(iocToReporterEdges, iocToReporterEdge)
		malwareToReporterEdges = append(malwareToReporterEdges, malwareToReporterEdge)
		reporterToIocEdges = append(reporterToIocEdges, reporterToIocEdge)
		reporterToMalwareEdges = append(reporterToMalwareEdges, reporterToMalwareEdge)
		reporterToThreatSourceEdges = append(reporterToThreatSourceEdges, reporterToThreatSourceEdge)
	}

	var vertices = make(map[string]interface{})
	vertices["ioc"] = iocVertices
	vertices["malware"] = malwareVertices
	vertices["reporter"] = reporterVertices
	vertices["threatSource"] = threatSourceVertices

	var edges = make(map[string][]Edge)
	// ioc
	edges["associates_with"] = iocToMalwareEdges
	edges["came_from"] = iocToThreatSourceEdges
	edges["found_by"] = iocToReporterEdges
	// malware
	edges["associates_with"] = malwareToIocEdges
	edges["reported_by"] = malwareToReporterEdges
	// reporter
	edges["found"] = reporterToIocEdges
	edges["identified"] = reporterToMalwareEdges
	edges["reported_to"] = reporterToThreatSourceEdges
	//threat source
	edges["shared"] = threatSourceToIocEdges

	return vertices, edges
}
