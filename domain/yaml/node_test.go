package yaml_test

import (
	dyaml "github.com/gidoichi/yaml-path/domain/yaml"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	yamlv3 "gopkg.in/yaml.v3"
)

var _ = Describe("Node", func() {
	var (
		node dyaml.Node
	)

	Describe("FindChildValueByKey()", func() {
		Context("called from node type is not mapping", func() {
			BeforeEach(func() {
				node = dyaml.Node{
					Kind: yamlv3.ScalarNode,
				}
			})
			It("should return empty string", func() {
				value := node.FindChildValueByKey("name")

				Expect(value).To(BeEmpty())
			})
		})

		Context("called from node type is mapping", func() {
			Context("with given key", func() {
				BeforeEach(func() {
					node = dyaml.Node{
						Kind: yamlv3.MappingNode,
						Content: []*yamlv3.Node{
							{
								Kind:  yamlv3.ScalarNode,
								Value: "name",
							},
							{
								Kind:  yamlv3.ScalarNode,
								Value: "some",
							},
						},
					}
				})
				It("should return corresponding value", func() {
					value := node.FindChildValueByKey("name")

					Expect(value).To(Equal("some"))
				})
			})

			Context("without given key", func() {
				BeforeEach(func() {
					node = dyaml.Node{
						Kind: yamlv3.MappingNode,
						Content: []*yamlv3.Node{
							{
								Kind:  yamlv3.ScalarNode,
								Value: "dummy",
							},
							{
								Kind:  yamlv3.ScalarNode,
								Value: "other",
							},
						},
					}
				})
				It("should return empty string", func() {
					value := node.FindChildValueByKey("name")

					Expect(value).To(BeEmpty())
				})
			})
		})
	})

	Describe("FindSequenceSelectionByMappingKey()", func() {
		Context("called from node type is not sequence", func() {
			BeforeEach(func() {
				node = dyaml.Node{
					Kind: yamlv3.ScalarNode,
				}
			})
			It("should return empty string", func() {
				value := node.FindSequenceSelectionByMappingKey(0, "name")

				Expect(value).To(BeEmpty())
			})
		})

		Context("called from node type is sequence", func() {
			Context("with a mapping having given key", func() {
				BeforeEach(func() {
					node = dyaml.Node{
						Kind: yamlv3.SequenceNode,
						Content: []*yamlv3.Node{
							{
								Kind: yamlv3.MappingNode,
								Content: []*yamlv3.Node{
									{
										Kind:  yamlv3.ScalarNode,
										Value: "name",
									},
									{
										Kind:  yamlv3.ScalarNode,
										Value: "some",
									},
								},
							},
						},
					}
				})
				It("should return corresponding value", func() {
					value := node.FindSequenceSelectionByMappingKey(0, "name")

					Expect(value).To(Equal("some"))
				})
			})

			Context("with mappings not having given key", func() {
				BeforeEach(func() {
					node = dyaml.Node{
						Kind: yamlv3.SequenceNode,
						Content: []*yamlv3.Node{
							{
								Kind: yamlv3.MappingNode,
								Content: []*yamlv3.Node{
									{
										Kind:  yamlv3.ScalarNode,
										Value: "dummy",
									},
									{
										Kind:  yamlv3.ScalarNode,
										Value: "other",
									},
								},
							},
						},
					}
				})
				It("should return empty string", func() {
					value := node.FindSequenceSelectionByMappingKey(0, "name")

					Expect(value).To(BeEmpty())
				})
			})

			Context("with 2 or more mappings having given key", func() {
				BeforeEach(func() {
					node = dyaml.Node{
						Kind: yamlv3.SequenceNode,
						Content: []*yamlv3.Node{
							{
								Kind: yamlv3.MappingNode,
								Content: []*yamlv3.Node{
									{
										Kind:  yamlv3.ScalarNode,
										Value: "name",
									},
									{
										Kind:  yamlv3.ScalarNode,
										Value: "some",
									},
								},
							},
							{
								Kind: yamlv3.MappingNode,
								Content: []*yamlv3.Node{
									{
										Kind:  yamlv3.ScalarNode,
										Value: "name",
									},
									{
										Kind:  yamlv3.ScalarNode,
										Value: "other",
									},
								},
							},
						},
					}
				})
				It("should return empty string", func() {
					value := node.FindSequenceSelectionByMappingKey(0, "name")

					Expect(value).To(BeEmpty())
				})
			})
		})
	})
})
