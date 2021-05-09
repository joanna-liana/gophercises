import React from 'react';
import ActionButton from './ActionButton';
import './StartScreen.scss';

function StartScreen({
  startStory,
  isLoading,
}: {
  startStory: () => void;
  isLoading: boolean;
}) {
  return (
    <div className="StartScreen">
      <header className="StartScreen-header">
        <h1>Choose your own adventure</h1>
        <ActionButton
          className="StartScreen-button"
          isActionInProgress={isLoading}
          actionHandler={startStory}
          label="Begin"
        />
      </header>
    </div>
  );
}

export default StartScreen;
