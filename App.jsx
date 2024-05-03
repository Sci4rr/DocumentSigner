function memoize(fn) {
  const cache = {};
  return function(...args) {
    const key = JSON.stringify(args);
    if (!cache[key]) {
      cache[key] = fn.apply(this, args);
    }
    return cache[key];
  };
}

import React from 'react';

const DocumentViewer = React.memo(function DocumentViewer(props) {
  // Component code here
});

import React, { useMemo } from 'react';

const MyComponent = (props) => {
  const memoizedValue = useMemo(() => computeExpensiveValue(props.dependency), [props.dependency]);
  return <div>{memoizedValue}</div>;
};