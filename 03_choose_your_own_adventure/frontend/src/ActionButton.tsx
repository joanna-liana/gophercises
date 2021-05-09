import React from 'react';
import './ActionButton.scss';
import Loader from 'react-loader-spinner';

interface Props {
  actionHandler: () => void;
  isActionInProgress: boolean;
  label: string;
}

function ActionButton({
  actionHandler,
  isActionInProgress,
  label,
  ...props
}: Props & React.HTMLAttributes<HTMLButtonElement>) {
  return (
    <button
      {...props}
      className={`${props.className} ActionButton`}
      onClick={actionHandler}
      disabled={isActionInProgress}
    >
      {isActionInProgress ? (
        <Loader type="Oval" width="30px" height="30px" />
      ) : (
        label
      )}
    </button>
  );
}

export default ActionButton;
