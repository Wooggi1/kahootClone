<script lang="ts">
  import Button from "../../lib/Button.svelte";
  import QuizCard from "../../lib/QuizCard.svelte";
  import type { Quiz } from "../../model/quiz";

  let quizzes: Quiz[] = [];

  async function createQuiz(): Promise<void> {
    const quizData = {
      name: "Teste",
      subjects: ["literature"] 
    };

    const response = await fetch("http://localhost:3000/api/quiz/create", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(quizData)
    });

    if (!response.ok) {
      alert("Failed to create quiz!");
      return;
    }

    const newQuiz = await response.json();
    quizzes = [...quizzes, newQuiz];
    console.log(newQuiz)
  }
</script>

<div class="p-8">
  <h2 class="text-4xl font-bold">Your quizzes</h2>
  <div class="flex flex-col gap-2 mt-4">
    {#each quizzes as quiz(quiz.id)}
      <QuizCard on:host {quiz} />
    {/each}
  </div>
  <Button on:click={createQuiz}>
    create quiz
  </Button>
</div>