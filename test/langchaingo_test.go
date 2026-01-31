package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/memory"
	"github.com/tmc/langchaingo/prompts"
)

var (
	llm *openai.LLM

	ctx           = context.Background()
	apiKey string = "79cfc89b64ea45f8b5ab83368844e32a.YgAMAzBH2uZPeA5N"
)

func init() {
	var err error
	llm, err = openai.New(
		openai.WithToken(apiKey),
		openai.WithBaseURL("https://open.bigmodel.cn/api/paas/v4"),
		openai.WithModel("glm-4"),
	)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// memory 记忆
func TestMemory(t *testing.T) {
	ctx := context.Background()
	message := "What is my name?"

	wb := memory.NewConversationBuffer(memory.WithAIPrefix("ni hao"))
	wb.ChatHistory.AddUserMessage(ctx, "Hi,My name is mengfanbing")
	wb.ChatHistory.AddAIMessage(ctx, "Can i help you?")

	response, err := chains.Run(ctx, chains.NewConversation(llm, wb), message)
	if err != nil {
		t.Fatalf("Failed to chat: %v", err)
	}
	fmt.Println(response)
}

func TestChains(t *testing.T) {
	ctx := context.Background()

	prompt := "What is the best name to describe acompany that makes {{.product}}?"

	chainOne := chains.NewLLMChain(llm, prompts.NewPromptTemplate(prompt, []string{"product"}))

	product := "Queen Size Sheet Set"
	response, err := chains.Run(ctx, chainOne, product)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Println("response:", response)

}

func TestSimpleSequentialChains(t *testing.T) {
	product := "加大号床单套装"

	ctx := context.Background()

	firstPrompt := "对于一家生产{{.product}}的公司，最好的名称是什么?"
	chainOne := chains.NewLLMChain(llm, prompts.NewPromptTemplate(firstPrompt, []string{"product"}))

	// response, err := chains.Run(ctx, chainOne, product)
	// if err != nil {
	// 	fmt.Println("err:", err)
	// 	return
	// }
	// fmt.Println("response:", response)

	secondPrompt := "为以下公司写一个20个字的描述:{{.companyName}}"
	chainTwo := chains.NewLLMChain(llm, prompts.NewPromptTemplate(secondPrompt, []string{"companyName"}))

	overallSimpleChain, err := chains.NewSimpleSequentialChain([]chains.Chain{chainOne, chainTwo})
	if err != nil {
		fmt.Println("err:", err)
		return
	}

	response, err := chains.Run(ctx, overallSimpleChain, product)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Println("response:", response)
}

// QA
// Q: 多个链不能包含相同的输入参数
// (err: invalid input values: missing key in input values: {相同的输入参数key})
func TestSequentialChains(t *testing.T) {
	review := "这个产品非常好用，质量很棒！"

	// firstPrompt := "翻译下面的评论为英语:{{.review}}"
	// chainOne := chains.NewLLMChain(llm, prompts.NewPromptTemplate(firstPrompt, []string{"review"}))
	// chainOne.OutputKey = "englishReview"

	secondPrompt := "用一句话总结下面的评论:{{.review}}"
	chainTwo := chains.NewLLMChain(llm, prompts.NewPromptTemplate(secondPrompt, []string{"review"}))
	chainTwo.OutputKey = "summary"

	thirdPrompt := "下面的评论是什么语言:{{.review}}"
	chainThird := chains.NewLLMChain(llm, prompts.NewPromptTemplate(thirdPrompt, []string{"review"}))
	chainThird.OutputKey = "language"

	// fourthPrompt := "使用特定的语言对下面的总结写一个回复:\n\n总结:{{.summary}}"
	fourthPrompt := "使用特定的语言对下面的总结写一个回复:\n\n总结:{{.summary}}\n\n语言:{{.language}}"
	chainFour := chains.NewLLMChain(llm, prompts.NewPromptTemplate(fourthPrompt, []string{"summary", "language"}))
	chainFour.OutputKey = "response"

	overallChain, err := chains.NewSequentialChain(
		[]chains.Chain{chainTwo, chainThird, chainFour},
		[]string{"review"},
		[]string{"response"})
	if err != nil {
		fmt.Println("err:", err)
		return
	}

	outputs, err := chains.Call(ctx, overallChain, map[string]any{"review": review})
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	// fmt.Println("outputs:", outputs)
	if response, ok := outputs["response"]; ok {
		fmt.Println("response:", response)
	}
}

func TestSequentialChain(t *testing.T) {
	ctx := context.Background()
	t.Parallel()

	chain1 := chains.NewLLMChain(
		llm,
		prompts.NewPromptTemplate("Write a story titled {{.title}} set in the year {{.year}}", []string{"title", "year"}),
	)
	chain1.OutputKey = "story"
	chain2 := chains.NewLLMChain(llm, prompts.NewPromptTemplate("Review this story: {{.story}}", []string{"story"}))
	chain2.OutputKey = "review"
	chain3 := chains.NewLLMChain(
		llm,
		prompts.NewPromptTemplate("Tell me if this review is legit: {{.review}}", []string{"review"}),
	)
	chain3.OutputKey = "result"

	chainss := []chains.Chain{chain1, chain2, chain3}

	seqChain, err := chains.NewSequentialChain(chainss, []string{"title", "year"}, []string{"result"})
	if err != nil {
		fmt.Println("chains.NewSequentialChain error:", err)
		return
	}

	res, err := chains.Call(ctx, seqChain, map[string]any{"title": "Chicken Takeover", "year": 3000})
	if err != nil {
		fmt.Println("chains.Call error:", err)
		return
	}

	fmt.Println("result:", res["result"])
}

// langchaingo不支持多提示词链?
func TestMultiPromptChain(t *testing.T) {

}
