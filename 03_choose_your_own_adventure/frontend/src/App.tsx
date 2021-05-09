import React, { useEffect, useState } from 'react';
import useFetch from 'use-http';
import StartScreen from './StartScreen';
import StoryCard from './StoryCard';
import './App.scss';

enum StoryState {
  START = 'START',
  IN_PROGRESS = 'IN_PROGRESS',
}

export interface StoryData {
  Title: string;
  Story: string[];
  Options: { Text: string; Arc: string }[];
}

function App() {
  const [storyState, setStoryState] = useState<StoryState>(StoryState.START);
  const [storyData, setStoryData] = useState<StoryData | null>(null);
  // TODO: error handling
  const { get, response, loading } = useFetch('http://localhost:8080');

  useEffect(() => {
    async function initialiseStory() {
      const story = await get('/start');
      if (response.ok) setStoryData(story);
    }

    initialiseStory();
  }, []);

  async function chooseNextStoryArc(option: string) {
    const story = await get(`/next/${option}`);
    if (response.ok) setStoryData(story);
  }

  async function startStory() {
    setStoryState(StoryState.IN_PROGRESS);
  }

  return (
    <div className="App-contentWrapper">
      {storyState === StoryState.START && (
        <StartScreen startStory={startStory} isLoading={loading} />
      )}

      {storyState === StoryState.IN_PROGRESS && storyData && (
        <StoryCard
          story={storyData}
          chooseNext={chooseNextStoryArc}
          isLoading={loading}
        />
      )}
    </div>
  );
}

export default App;
