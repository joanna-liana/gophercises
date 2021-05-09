import React from 'react';
import ActionButton from './ActionButton';
import { StoryData } from './App';
import './StoryCard.scss';

// TODO: show each part on click

function StoryCard({
  story,
  chooseNext,
  isLoading,
}: {
  story: StoryData;
  chooseNext: (option: string) => Promise<void>;
  isLoading: boolean;
}) {
  return (
    <div className="StoryCard">
      <div>{story.Story}</div>
      {story.Options.map((option, i) => (
        <ActionButton
          isActionInProgress={isLoading}
          label={option.Text}
          actionHandler={() => chooseNext(option.Arc)}
          key={i}
        />
      ))}
    </div>
  );
}

export default StoryCard;
