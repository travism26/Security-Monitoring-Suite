export const kafkaWrapper = {
  getProducer: jest.fn().mockImplementation((producerName) => {
    console.log("Mocked Producer found for:", producerName);
    return {
      publish: jest.fn().mockImplementation((data) => {
        console.log("Mocked Producer publish:", data);
        return Promise.resolve();
      }),
    };
  }),
  getConsumer: jest.fn().mockImplementation(() => {
    console.log("Mocked Consumer found");
    return {
      subscribe: jest.fn().mockResolvedValue(null),
      run: jest.fn().mockResolvedValue(null),
      connect: jest.fn().mockResolvedValue(null),
      disconnect: jest.fn().mockResolvedValue(null),
    };
  }),
  isInitialized: jest.fn().mockReturnValue(true),
};
